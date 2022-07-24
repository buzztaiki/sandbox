use std::cell::{Cell, RefCell};
use std::collections::HashMap;
use std::io;
use std::rc::Rc;

use crate::disk::{DiskManager, Page, PageId, PAGE_SIZE};

#[derive(Debug, thiserror::Error)]
pub enum Error {
    #[error(transparent)]
    Io(#[from] io::Error),
    #[error("no free buffer")]
    NoFreeBuffer,
}

#[derive(Debug, Clone, Copy, Default)]
pub struct BufferId(u64);

pub struct Buffer {
    pub page_id: PageId,
    pub page: RefCell<Page>,
    pub is_dirty: Cell<bool>,
}

impl Default for Buffer {
    fn default() -> Self {
        Buffer {
            page_id: PageId::default(),
            page: RefCell::new([0; PAGE_SIZE]),
            is_dirty: Cell::new(false),
        }
    }
}

#[derive(Default)]
pub struct Frame {
    usage_count: u64,
    buffer: Rc<Buffer>,
}

pub struct BufferPool {
    buffers: Vec<Frame>,
    next_victim_id: BufferId,
}

impl BufferPool {
    pub fn new(size: usize) -> Self {
        let mut buffers = Vec::new();
        buffers.resize_with(size, Frame::default);
        BufferPool {
            buffers,
            next_victim_id: BufferId::default(),
        }
    }

    // 捨てるバッファの BufferId を返す。
    pub fn evict(&mut self) -> Option<BufferId> {
        let pool_size = self.size();
        let mut consecutive_pinned = 0;

        // Clock-sweep algorithm
        loop {
            let next_victim_id = self.next_victim_id;
            let frame = &mut self[next_victim_id];
            if frame.usage_count == 0 {
                break Some(self.next_victim_id);
            }

            // get_mut できるという事は参照が存在しないということ。
            // 利用の管理だけなら、Rc::get_mut だけで十分だが、利用回数が少ないものを優先で返したい意図があるため usage_count を使ってる。はず。
            if Rc::get_mut(&mut frame.buffer).is_some() {
                frame.usage_count -= 1;
                consecutive_pinned = 0;
            } else {
                consecutive_pinned += 1;
                if consecutive_pinned >= pool_size {
                    break None;
                }
            }

            self.next_victim_id = self.rotate_id(self.next_victim_id);
        }
    }

    pub fn size(&self) -> usize {
        self.buffers.len()
    }

    fn rotate_id(&self, id: BufferId) -> BufferId {
        BufferId((id.0 + 1) % self.size() as u64)
    }
}

impl std::ops::Index<BufferId> for BufferPool {
    type Output = Frame;

    fn index(&self, index: BufferId) -> &Frame {
        &self.buffers[index.0 as usize]
    }
}

impl std::ops::IndexMut<BufferId> for BufferPool {
    fn index_mut(&mut self, index: BufferId) -> &mut Self::Output {
        &mut self.buffers[index.0 as usize]
    }
}

pub struct BufferPoolManager {
    disk: DiskManager,
    pool: BufferPool,
    page_table: HashMap<PageId, BufferId>,
}

impl BufferPoolManager {
    pub fn new(disk: DiskManager, pool: BufferPool) -> Self {
        let page_table = HashMap::new();
        BufferPoolManager {
            disk,
            pool,
            page_table,
        }
    }

    pub fn create_page(&mut self) -> Result<Rc<Buffer>, Error> {
        self.assign_page(None)
    }

    pub fn fetch_page(&mut self, page_id: PageId) -> Result<Rc<Buffer>, Error> {
        if let Some(&buffer_id) = self.page_table.get(&page_id) {
            let frame = &mut self.pool[buffer_id];
            frame.usage_count += 1;
            Ok(Rc::clone(&frame.buffer))
        } else {
            self.assign_page(Some(page_id))
        }
    }

    fn assign_page(&mut self, page_id: Option<PageId>) -> Result<Rc<Buffer>, Error> {
        let buffer_id = self.pool.evict().ok_or(Error::NoFreeBuffer)?;
        let frame = &mut self.pool[buffer_id];

        // evict できているという事は所有者がいないという事
        let buffer = Rc::get_mut(&mut frame.buffer).unwrap();
        let old_page_id = buffer.page_id;
        if buffer.is_dirty.get() {
            self.disk.write_page_data(old_page_id, buffer.page.get_mut())?;
        }
        self.page_table.remove(&old_page_id);

        let page_id = match page_id {
            Some(page_id) => {
                self.disk.read_page_data(page_id, buffer.page.get_mut())?;
                buffer.is_dirty.set(false);
                page_id
            },
            None => {
                *buffer = Buffer::default();
                buffer.is_dirty.set(true);
                self.disk.allocate_page()
            },
        };
        buffer.page_id = page_id;
        self.page_table.insert(page_id, buffer_id);

        frame.usage_count = 1;
        Ok(Rc::clone(&frame.buffer))
    }

    pub fn flush(&mut self) -> Result<(), Error> {
        for (&page_id, &buffer_id) in self.page_table.iter() {
            let frame = &self.pool[buffer_id];
            let mut page = frame.buffer.page.borrow_mut();
            self.disk.write_page_data(page_id, page.as_mut())?;
            frame.buffer.is_dirty.set(false);
        }
        self.disk.sync()?;
        Ok(())
    }
}

#[cfg(test)]
mod tests {
    use super::*;
    use tempfile::tempfile;

    #[test]
    fn test() -> anyhow::Result<()> {
        let mut hello = Vec::with_capacity(PAGE_SIZE);
        hello.extend_from_slice(b"hello");
        hello.resize(PAGE_SIZE, 0);
        let mut world = Vec::with_capacity(PAGE_SIZE);
        world.extend_from_slice(b"world");
        world.resize(PAGE_SIZE, 0);

        let disk = DiskManager::new(tempfile().unwrap()).unwrap();
        let pool = BufferPool::new(1);
        let mut bufmgr = BufferPoolManager::new(disk, pool);
        let page1_id = {
            let buffer = bufmgr.create_page().unwrap();
            assert!(bufmgr.create_page().is_err());
            let mut page = buffer.page.borrow_mut();
            page.copy_from_slice(&hello);
            buffer.is_dirty.set(true);
            buffer.page_id
        };
        {
            let buffer = bufmgr.fetch_page(page1_id).unwrap();
            let page = buffer.page.borrow();
            assert_eq!(&hello, page.as_ref());
        }
        let page2_id = {
            let buffer = bufmgr.create_page().unwrap();
            let mut page = buffer.page.borrow_mut();
            page.copy_from_slice(&world);
            buffer.is_dirty.set(true);
            buffer.page_id
        };

        bufmgr.flush()?;
        {
            assert_ne!(page1_id, page2_id);

            let mut buf = vec![0; PAGE_SIZE];
            bufmgr.disk.read_page_data(page1_id, &mut buf).unwrap();
            assert_eq!(hello, buf);
            bufmgr.disk.read_page_data(page2_id, &mut buf).unwrap();
            assert_eq!(world, buf);

        }

        {
            let buffer = bufmgr.fetch_page(page1_id).unwrap();
            let page = buffer.page.borrow();
            // FAIL: world に入れかわってる
            assert!(&hello == page.as_ref());
            // assert_eq!(&hello, page.as_ref());
        }
        {
            let buffer = bufmgr.fetch_page(page2_id).unwrap();
            let page = buffer.page.borrow();
            assert_eq!(&world, page.as_ref());
        }

        Ok(())
    }
}
