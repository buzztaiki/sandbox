use anyhow::Result;

use relly::btree::{BTree, SearchMode};
use relly::buffer::{BufferPool, BufferPoolManager};
use relly::disk::{DiskManager, PageId};

fn main() -> Result<()> {
    let disk = DiskManager::open("test.btr")?;
    let pool = BufferPool::new(10);
    let mut bufmgr = BufferPoolManager::new(disk, pool);

    // TODO: btree.search で panic する。
    // btree-create で btree をそのまま再利用したら panic しなかったから、ディスクへの書き込みか読み込みのどっちかがおかしい気がする。
    // 元の relly にテストが付いてたはずだから、それ使って確認してみるのがよい？
    let btree = BTree::new(PageId(0));
    let mut iter = btree.search(&mut bufmgr, SearchMode::Key(b"Hyogo".to_vec()))?;
    let (key, value) = iter.next(&mut bufmgr)?.unwrap();
    println!("{} = {}", String::from_utf8(key)?, String::from_utf8(value)?);
    Ok(())
}
