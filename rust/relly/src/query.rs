use anyhow::Result;

use crate::btree::{BTree, SearchMode};
use crate::buffer::BufferPoolManager;
use crate::disk::PageId;
use crate::{btree, tuple};

pub type Tuple = Vec<Vec<u8>>;
pub type TupleSlice<'a> = &'a [Vec<u8>];

pub enum TupleSearchMode<'a> {
    Start,
    Key(&'a [&'a [u8]]),
}

impl<'a> TupleSearchMode<'a> {
    fn encode(&self) -> SearchMode {
        match self {
            TupleSearchMode::Start => SearchMode::Start,
            TupleSearchMode::Key(tuple) => {
                let mut key = vec![];
                tuple::encode(tuple.iter(), &mut key);
                SearchMode::Key(key)
            }
        }
    }
}

pub trait Executor {
    fn next(&mut self, bufmgr: &mut BufferPoolManager) -> Result<Option<Tuple>>;
}

pub type BoxExecutor<'a> = Box<dyn Executor + 'a>;

pub trait PlanNode {
    fn start(&self, bufmgr: &mut BufferPoolManager) -> Result<BoxExecutor>;
}

pub struct SeqScan<'a> {
    pub table_meta_page_id: PageId,
    pub search_mode: TupleSearchMode<'a>,
    pub while_cond: &'a dyn Fn(TupleSlice) -> bool,
}

impl<'a> PlanNode for SeqScan<'a> {
    fn start(&self, bufmgr: &mut BufferPoolManager) -> Result<BoxExecutor> {
        let btree = BTree::new(self.table_meta_page_id);
        let table_iter = btree.search(bufmgr, self.search_mode.encode())?;
        Ok(Box::new(ExecSeqScan {
            table_iter,
            while_cond: self.while_cond,
        }))
    }
}

pub struct ExecSeqScan<'a> {
    table_iter: btree::Iter,
    while_cond: &'a dyn Fn(TupleSlice) -> bool,
}

impl<'a> Executor for ExecSeqScan<'a> {
    fn next(&mut self, bufmgr: &mut BufferPoolManager) -> Result<Option<Tuple>> {
        let (key_bytes, tuple_bytes) = match self.table_iter.next(bufmgr)? {
            Some(pair) => pair,
            None => return Ok(None),
        };

        let mut key = vec![];
        tuple::decode(&key_bytes, &mut key);
        if !(self.while_cond)(&key) {
            return Ok(None);
        }
        let mut tuple = key;
        tuple::decode(&tuple_bytes, &mut tuple);

        Ok(Some(tuple))
    }
}

pub struct Filter<'a> {
    pub inner_plan: &'a dyn PlanNode,
    pub cond: &'a dyn Fn(TupleSlice) -> bool,
}

impl<'a> PlanNode for Filter<'a> {
    fn start(&self, bufmgr: &mut BufferPoolManager) -> Result<BoxExecutor> {
        let inner_iter = self.inner_plan.start(bufmgr)?;
        Ok(Box::new(ExecFilter {
            inner_iter,
            cond: self.cond,
        }))
    }
}

pub struct ExecFilter<'a> {
    inner_iter: BoxExecutor<'a>,
    cond: &'a dyn Fn(TupleSlice) -> bool,
}

impl<'a> Executor for ExecFilter<'a> {
    fn next(&mut self, bufmgr: &mut BufferPoolManager) -> Result<Option<Tuple>> {
        loop {
            match self.inner_iter.next(bufmgr)? {
                Some(tuple) => {
                    if (self.cond)(&tuple) {
                        return Ok(Some(tuple));
                    }
                }
                None => return Ok(None),
            }
        }
    }
}

pub struct IndexScan<'a> {
    pub table_meta_page_id: PageId,
    pub index_meta_page_id: PageId,
    pub search_mode: TupleSearchMode<'a>,
    pub while_cond: &'a dyn Fn(TupleSlice) -> bool,
}

impl<'a> PlanNode for IndexScan<'a> {
    fn start(&self, bufmgr: &mut BufferPoolManager) -> Result<BoxExecutor> {
        let table_btree = BTree::new(self.table_meta_page_id);
        let index_btree = BTree::new(self.index_meta_page_id);
        let index_iter = index_btree.search(bufmgr, self.search_mode.encode())?;
        Ok(Box::new(ExecIndexScan {
            table_btree,
            index_iter,
            while_cond: self.while_cond,
        }))
    }
}

pub struct ExecIndexScan<'a> {
    table_btree: BTree,
    index_iter: btree::Iter,
    while_cond: &'a dyn Fn(TupleSlice) -> bool,
}

impl<'a> Executor for ExecIndexScan<'a> {
    fn next(&mut self, bufmgr: &mut BufferPoolManager) -> Result<Option<Tuple>> {
        let (skey_bytes, pkey_bytes) = match self.index_iter.next(bufmgr)? {
            Some(pair) => pair,
            None => return Ok(None),
        };

        let mut skey = vec![];
        tuple::decode(&skey_bytes, &mut skey);
        if !(self.while_cond)(&skey) {
            return Ok(None);
        }

        let mut table_iter = self
            .table_btree
            .search(bufmgr, SearchMode::Key(pkey_bytes))?;
        let (pkey_bytes, tuple_bytes) = table_iter.next(bufmgr)?.unwrap();
        let mut tuple = vec![];
        tuple::decode(&pkey_bytes, &mut tuple);
        tuple::decode(&tuple_bytes, &mut tuple);
        Ok(Some(tuple))
    }
}
