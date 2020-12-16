use std::sync::Arc;
use std::sync::atomic::{AtomicU16, Ordering, AtomicU64};
use std::thread;

const TOTAL_SIZE:usize = 100 * 100;
const NTHREAD:usize = 6;

fn main() {
    let data : Vec<i32> = (1..(TOTAL_SIZE+1) as i32).collect(); //初始化一个数据从1到n数组
    let arc_data = Arc::new(data); //data 的所有权转给了 ar_data
    let result  = Arc::new(AtomicU64::new(0)); //收集结果的数组(原子操作)
    let mut thread_handlers = vec![]; // 用于收集线程句柄
    for i in 0..NTHREAD {
        // clone Arc 准备move到线程中，只增加引用计数，不会深拷贝内部数据
        let test_data = arc_data.clone();
        let res = result.clone();
        thread_handlers.push(
            thread::spawn(move || {
                let id = i;
                //找到自己的分区
                let chunk_size = TOTAL_SIZE / NTHREAD + 1;
                let start = id * chunk_size;
                let end = std::cmp::min(start + chunk_size, TOTAL_SIZE);
                //进行求和运算
                let mut sum = 0;
                for  i in start..end  {
                    sum += test_data[i];
                }
                //原子操作
                res.fetch_add(sum as u64, Ordering::SeqCst);
                println!("id={}, sum={}", id, sum );
            }
            ));
    }
    // 等所有的线程执行完
    for th in thread_handlers {
        th.join().expect("The sender thread panic!!!");
    }
    // 输出结果
    println!("result = {}",result.load(Ordering::SeqCst));
}