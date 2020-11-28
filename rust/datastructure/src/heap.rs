use std::process::id;

#[derive(Debug)]
struct Heap {
    data: Vec<Option<i32>>,
    capacity: usize,
    count: i32,
}

impl Heap {
    pub fn new(capacity: usize) -> Self {
        Heap {
            data: vec![None; capacity],
            capacity: capacity,
            count: 0
        }
    }

    pub fn insert(&mut self, x: i32) -> bool {
        if self.capacity as i32 == self.count { false }

        self.data[self.count as usize] = Some(x);
        if self.count == 0 {
            self.count += 1;
            true
        }

        let mut idx = self.count as usize;

        let mut parent_idx = ((idx - 1) >> 1) as usize;
        while parent_idx > 0 && self.data[idx] > self.data[parent_idx] {
            self.swap(idx, parent_idx);
            idx = parent_idx;
            parent_idx = ((idx - 1) >> 1) as usize;
        }
        self.count += 1;
        true
    }

    pub fn remove_max(&mut self) -> Option<i32> {
        if self.count == 0 { None }
        let max_value = self.data[0];

        self.data[0] = self.data[(self.count - 1) as usize];
        self.data[(self.count - 1) as usize] = None;

        self.heapify();
        self.count -= 1;
        max_value
    }

    pub fn heapify(&mut self) {
        let mut idx = 0usize;
        loop {
            let mut max_pos = idx;
            if (2 * idx + 1) as i32 <= self.count && self.data[idx] < self.data[2 * idx + 1] { max_pos = 2 * idx + 1 }
            if (2 * idx + 2) as i32 <= self.count && self.data[max_pos] < self.data[2 * idx + 2] { max_pos = 2 * idx + 2 }

            if max_pos == idx { break }
            self.swap(idx, max_pos);
            idx = max_pos
        }
    }

    pub fn swap(&mut self, idx: usize, parent_idx: usize) {
        let tmp = self.data[parent_idx];
        self.data[parent_idx] = self.data[idx];
        self.data[idx] = tmp;
    }
}



