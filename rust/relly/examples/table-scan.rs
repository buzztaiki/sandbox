use anyhow::Result;

use relly::buffer::{BufferPool, BufferPoolManager};
use relly::disk::{DiskManager, PageId};
use relly::query::{Filter, PlanNode, SeqScan, TupleSearchMode};
use relly::tuple;

// SELECT * WHERE last_name = 'Smith'
// without index
fn main() -> Result<()> {
    let disk = DiskManager::open("table.rly")?;
    let pool = BufferPool::new(10);
    let mut bufmgr = BufferPoolManager::new(disk, pool);

    let plan = Filter {
        cond: &|record| record[2].as_slice() == b"Smith",
        inner_plan: &SeqScan {
            table_meta_page_id: PageId(0),
            search_mode: TupleSearchMode::Start,
            while_cond: &|_| true,
        },
    };
    let mut exec = plan.start(&mut bufmgr)?;

    while let Some(record) = exec.next(&mut bufmgr)? {
        println!("{:?}", tuple::Pretty(&record));
    }
    Ok(())
}
