
fn bubble_sort(mut nums: Vec<i32>) -> Vec<i32> {
    if nums.is_empty() { vec![]; }

    let n = nums.len();

    for i in 0..n {
        let mut swap = false;
        for j in 0..n-i-1 {
            if nums[j] > nums[j+1] {
                swap = true;
                let tmp = nums[j];
                nums[j] = nums[j + 1];
                nums[j+1] = tmp;
            }
        }
        if !swap { break;}
    }
    nums
}