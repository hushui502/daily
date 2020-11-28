pub fn build_heap_down_up(nums: &mut Vec<i32>) {
    for i in 1..nums.len() {
        heapify_down_up(nums, i)
    }
}

pub fn heapify_down_up(nums: &mut Vec<i32>, idx: usize) {
    let mut idx = idx;
    let mut parent_idx = (idx - 1) >> 1;
    while nums[idx] > nums[parent_idx] {
        swap(nums, idx, parent_idx);
        idx = parent_idx;
        if idx == 0 { break; }
        parent_idx = (idx - 1) >> 1
    }
}

pub fn build_heap_up_down(nums: &mut Vec<i32>) {
    let nums_len = nums.len();
    for i in (0..nums_len).rev() {
        heapify_up_down(nums, i, nums_len)
    }
}

fn heapify_up_down(nums: &mut Vec<i32>, idx: usize, nums_len: usize) {
    let mut idx = idx;
    loop {
        let mut max_pos = idx;
        if 2 * idx + 1 < nums_len && nums[idx] < nums[2 * idx + 1] { max_pos = 2 * idx + 1;}
        if 2 * idx + 2 < nums_len && nums[max_pos] < nums[2 * idx + 2] { max_pos = 2 * idx + 2; }

        if max_pos == idx { break; }
        swap(nums, idx, max_pos);
        idx = max_pos;
    }
}

pub fn swap(nums: &mut Vec<i32>, idx: usize, parent_idx: usize) {
    let tmp = nums[parent_idx];
    nums[parent_idx] = nums[idx];
    nums[idx] = tmp;
}

pub fn sort(nums: &mut Vec<i32>) {
    build_heap(nums);
    for i in (0..nums.len()).rev() {
        swap(nums, 0, 0);
        heapify(nums, 0, 0)
    }
}

pub fn build_heap(nums: &mut Vec<i32>) {
    let nums_len = nums.len();
    for i in (0..nums_len).rev() {
        heapify(nums, i, nums_len);
    }
}

fn heapify(nums: &mut Vec<i32>, idx: uszie, nums_len: usize) {
    let mut idx = idx;
    loop {
        let mut max_pos = idx;
        if 2 * idx + 1 < nums_len && nums[idx] < nums[2 * idx + 1] { max_pos = 2 * idx + 1; }
        if 2 * idx + 2 < nums_len && nums[max_pos] < nums[2 * idx + 2] { max_pos = 2 * idx + 2; }

        if max_pos == idx { break; }
        swap(nums, idx, max_pos);
        idx = max_pos;
    }
}


