use std::borrow::BorrowMut;
use std::cell::RefCell;
use std::collections::{HashMap, HashSet, VecDeque};
use std::fmt::format;
use std::process::id;
use std::rc::Rc;

mod case {
    use std::fmt;
    use std::fmt::Formatter;
    use std::io::Read;

    pub struct BuffReader<R> {
        inner: R,
        buf: Box<[u8]>,
        pos: usize,
        cap: usize,
    }

    impl<R> fmt::Debug for BuffReader<R>
    where
        R: fmt::Debug,
    {
        fn fmt(&self, f: &mut Formatter<'_>) -> fmt::Result {
            todo!()
        }
    }
}

// 一个名为 `my_mod` 的模块
mod my_mod {
    // 模块中的项默认具有私有的可见性
    fn private_function() {
        println!("called `my_mod::private_function()`");
    }

    // 使用 `pub` 修饰语来改变默认可见性。
    pub fn function() {
        println!("called `my_mod::function()`");
    }

    // 在同一模块中，项可以访问其它项，即使它是私有的。
    pub fn indirect_access() {
        print!("called `my_mod::indirect_access()`, that\n> ");
        private_function();
    }

    // 模块也可以嵌套
    pub mod nested {
        pub fn function() {
            println!("called `my_mod::nested::function()`");
        }

        #[allow(dead_code)]
        fn private_function() {
            println!("called `my_mod::nested::private_function()`");
        }

        // 使用 `pub(in path)` 语法定义的函数只在给定的路径中可见。
        // `path` 必须是父模块（parent module）或祖先模块（ancestor module）
        pub(in crate::my_mod) fn public_function_in_my_mod() {
            print!("called `my_mod::nested::public_function_in_my_mod()`, that\n > ");
            public_function_in_nested()
        }

        // 使用 `pub(self)` 语法定义的函数则只在当前模块中可见。
        pub(self) fn public_function_in_nested() {
            println!("called `my_mod::nested::public_function_in_nested");
        }

        // 使用 `pub(super)` 语法定义的函数只在父模块中可见。
        pub(super) fn public_function_in_super_mod() {
            println!("called my_mod::nested::public_function_in_super_mod");
        }
    }

    pub fn call_public_function_in_my_mod() {
        print!("called `my_mod::call_public_funcion_in_my_mod()`, that\n> ");
        nested::public_function_in_my_mod();
        print!("> ");
        nested::public_function_in_super_mod();
    }

    // `pub(crate)` 使得函数只在当前 crate 中可见
    pub(crate) fn public_function_in_crate() {
        println!("called `my_mod::public_function_in_crate()");
    }

    // 嵌套模块的可见性遵循相同的规则
    mod private_nested {
        #[allow(dead_code)]
        pub fn function() {
            println!("called `my_mod::private_nested::function()`");
        }
    }
}

fn function() {
    println!("called `function()`");
}

fn main() {
    // 模块机制消除了相同名字的项之间的歧义。
    function();
    my_mod::function();

    // 公有项，包括嵌套模块内的，都可以在父模块外部访问。
    my_mod::indirect_access();
    my_mod::nested::function();
    my_mod::call_public_function_in_my_mod();

    // pub(crate) 项可以在同一个 crate 中的任何地方访问
    my_mod::public_function_in_crate();

    // pub(in path) 项只能在指定的模块中访问
    // 报错！函数 `public_function_in_my_mod` 是私有的
    //my_mod::nested::public_function_in_my_mod();
    // 试一试 ^ 取消该行的注释

    // 模块的私有项不能直接访问，即便它是嵌套在公有模块内部的

    // 报错！`private_function` 是私有的
    //my_mod::private_function();
    // 试一试 ^ 取消此行注释

    // 报错！`private_function` 是私有的
    //my_mod::nested::private_function();
    // 试一试 ^ 取消此行的注释

    // Error! `private_nested` is a private module
    //my_mod::private_nested::function();
    // 试一试 ^ 取消此行的注释
}

pub fn move_zeroes(nums: &mut Vec<i32>) {
    let mut i = 0;
    let mut j = 0;
    while i < nums.len() {
        if nums[i] != 0 {
            nums.swap(i, j);
            j += 1;
        }
        i += 1;
    }
}

pub fn plus_one(mut digits: Vec<i32>) -> Vec<i32> {
    for i in (0..digits.len()).rev() {
        if digits[i] != 9 {
            digits[i] += 1;
            return digits;
        }
        digits[i] = 0;
        if i == 0 {
            digits.insert(0, 1);
        }
    }
    digits
}

pub fn remove_duplicates(nums: &mut Vec<i32>) -> i32 {
    if nums.len() == 0 {
        return 0;
    }
    let mut i = 0;
    for j in 1..nums.len() {
        if nums[i] != nums[j] {
            if j - i > 1 {
                nums[i + 1] = nums[j]
            }
            i += 1;
        }
    }

    (i + 1) as i32
}

struct MinStack {
    stack: Vec<i32>,
    min_stack: Vec<i32>,
}

impl MinStack {
    fn new() -> Self {
        MinStack {
            stack: Vec::new(),
            min_stack: Vec::new(),
        }
    }

    fn push(&mut self, x: i32) {
        self.stack.push(x);
        if self.min_stack.is_empty() || x <= self.min_stack[self.min_stack.len() - 1] {
            self.min_stack.push(x);
        }
    }

    fn pop(&mut self) {
        if self.stack.pop().unwrap() == *self.min_stack.last().unwrap() {
            self.min_stack.pop();
        }
    }

    fn top(&self) -> i32 {
        self.stack[self.stack.len() - 1]
    }

    fn get_min(&self) -> i32 {
        self.min_stack[self.min_stack.len() - 1]
    }
}

fn is_valid(s: String) -> bool {
    let mut stack = Vec::new();
    for c in s.chars() {
        match c {
            '(' | '{' | '[' => {
                stack.push(c);
            }
            ')' => {
                if stack.pop().unwrap() != '(' {
                    return false;
                }
            }
            '}' => {
                if stack.pop().unwrap() != '{' {
                    return false;
                }
            }
            ']' => {
                if stack.pop().unwrap() != '[' {
                    return false;
                }
            }
            _ => {}
        }
    }

    stack.is_empty()
}

pub fn max_sliding_window(nums: Vec<i32>, k: i32) -> Vec<i32> {
    if nums.len() == 0 || k == 1 {
        return nums;
    }

    let mut res: Vec<i32> = Vec::new();
    let mut deque: VecDeque<i32> = VecDeque::new();
    for i in 0..nums.len() {
        push(&mut deque, nums[i]);
        if (i as i32) > k - 1 {
            pop(&mut deque, nums[i - k as usize]);
            res.push(max(&deque));
        } else if (i as i32) == k - 1 {
            res.push(max(&deque));
        }
    }

    res
}

pub fn push(deque: &mut VecDeque<i32>, num: i32) {
    while !deque.is_empty() && *deque.back().unwrap() < num {
        deque.pop_back();
    }

    deque.push_back(num);
}

pub fn pop(deque: &mut VecDeque<i32>, n: i32) {
    if !deque.is_empty() && *deque.front().unwrap() == n {
        deque.pop_front();
    }
}

fn max(deque: &VecDeque<i32>) -> i32 {
    if deque.is_empty() {
        return 0;
    }

    *deque.front().unwrap()
}

pub fn two_sum(nums: Vec<i32>, target: i32) -> Vec<i32> {
    let mut map: HashMap<i32, i32> = HashMap::new();
    for i in 0..nums.len() {
        let complement = target - nums[i];
        if map.contains_key(&complement) {
            return vec![map[&complement], i as i32];
        }
        map.insert(nums[i], i as i32);
    }

    vec![]
}

fn is_anagram(s: String, t: String) -> bool {
    let mut map: HashMap<char, i32> = HashMap::new();
    for c in s.chars() {
        *map.entry(c).or_insert(0) += 1;
    }

    for c in t.chars() {
        if !map.contains_key(&c) {
            return false;
        }
        *map.entry(c).or_insert(0) -= 1;
        if *map.get(&c).unwrap() == 0 {
            map.remove(&c);
        }
    }

    map.is_empty()
}

pub fn group_anagram(strs: Vec<String>) -> Vec<Vec<String>> {
    let mut map: HashMap<String, Vec<String>> = HashMap::new();
    for s in strs {
        let mut chars: Vec<char> = s.chars().collect();
        chars.sort();
        let key = chars.iter().collect::<String>();
        map.entry(key).or_insert(Vec::new()).push(s);
    }

    map.values().map(|v| v.to_vec()).collect()
}

#[derive(Debug, PartialEq, Eq)]
pub struct ListNode {
    pub val: i32,
    pub next: Option<Box<ListNode>>,
}

impl ListNode {
    #[inline]
    fn new(val: i32) -> Self {
        ListNode { next: None, val }
    }

    // pub fn middle_node(head: Option<Box<ListNode>>) -> Option<Box<ListNode>> {
    //     let mut fast = head.as_ref().unwrap();
    //     let mut slow = head.as_ref();
    //
    //     if fast.next.is_none() { return head; }
    //
    //     while !fast.next.is_none() && !fast.next.as_ref().unwrap().next.is_none() {
    //         fast = fast.next.as_ref().unwrap().next.as_ref().unwrap();
    //         if fast.next.is_none() {
    //             return slow.unwrap().next.clone();
    //         }
    //         slow = slow.unwrap().next.as_ref();
    //
    //     }
    //     slow.unwrap().next.clone()
    // }

    // pub fn reverse_list(head: Option<Box<ListNode>>) -> Option<Box<ListNode>> {
    //     if head.is_none() {
    //         return None;
    //     }
    //     let mut prev: Option<Box<ListNode>> = None;
    //     let mut curr = head;
    //     while let Some(mut node) = curr {
    //         let next = node.next.take();
    //         node.next = prev.take();
    //         prev = Some(node);
    //         curr = next;
    //     }
    //
    //     prev
    // }
}

pub fn merge_two_lists(
    l1: Option<Box<ListNode>>,
    l2: Option<Box<ListNode>>,
) -> Option<Box<ListNode>> {
    match (l1, l2) {
        (Some(n1), Some(n2)) => {
            if n1.val < n2.val {
                Some(Box::new(ListNode {
                    val: n1.val,
                    next: merge_two_lists(n1.next, Some(n2)),
                }))
            } else {
                Some(Box::new(ListNode {
                    val: n2.val,
                    next: merge_two_lists(Some(n1), n2.next),
                }))
            }
        }
        (Some(n1), None) => Some(n1),
        (None, Some(n2)) => Some(n2),
        _ => None,
    }
}

pub fn remove_nth_from_end(head: Option<Box<ListNode>>, n: i32) -> Option<Box<ListNode>> {
    let mut dummy = Some(Box::new(ListNode { val: 0, next: head }));
    let mut cur = &mut dummy;
    let mut length = 0;

    while let Some(node) = cur.as_mut() {
        cur = &mut node.next;
        if let Some(_node) = cur {
            length += 1;
        }
    }

    let mut new_cur = dummy.as_mut();
    let idx = length - n;

    for _ in 0..idx {
        new_cur = new_cur.unwrap().next.as_mut();
    }

    let next = new_cur.as_mut().unwrap().next.as_mut().unwrap().next.take();
    new_cur.as_mut().unwrap().next = next;

    dummy.unwrap().next
}

#[derive(Debug, PartialEq, Eq)]
pub struct TreeNode {
    pub val: i32,
    pub left: Option<Rc<RefCell<TreeNode>>>,
    pub right: Option<Rc<RefCell<TreeNode>>>,
}

impl TreeNode {
    #[inline]
    pub fn new(val: i32) -> Self {
        TreeNode {
            val,
            left: None,
            right: None,
        }
    }

    pub fn preorder_traversal(root: Option<Rc<RefCell<TreeNode>>>) -> Vec<i32> {
        let mut res = Vec::new();

        if let Some(node) = root {
            res.push(node.borrow().val);
            res.append(&mut TreeNode::preorder_traversal(
                node.borrow().left.clone(),
            ));
            res.append(&mut TreeNode::preorder_traversal(
                node.borrow().right.clone(),
            ));
        }
        res
    }

    pub fn level_order(root: Option<Rc<RefCell<TreeNode>>>) -> Vec<Vec<i32>> {
        let mut levels: Vec<Vec<i32>> = Vec::new();
        if root.is_none() {
            return levels;
        }

        let mut deque: VecDeque<Option<Rc<RefCell<TreeNode>>>> = VecDeque::new();
        deque.push_back(root);

        while !deque.is_empty() {
            let mut current_level = vec![];
            let level_length = deque.len();
            for _ in 0..level_length {
                let node = deque.pop_front().unwrap();
                if let Some(node) = node {
                    current_level.push(node.borrow().val);
                    if node.borrow().left.is_some() {
                        deque.push_back(node.borrow().left.clone());
                    }
                    if node.borrow().right.is_some() {
                        deque.push_back(node.borrow().right.clone());
                    }
                }
            }
            levels.push(current_level);
        }

        levels
    }

    pub fn max_depth(root: Option<Rc<RefCell<TreeNode>>>) -> i32 {
        if root.is_none() {
            return 0;
        }

        let mut deque: VecDeque<Option<Rc<RefCell<TreeNode>>>> = VecDeque::new();
        deque.push_back(root);

        let mut depth = 0;
        while !deque.is_empty() {
            let level_length = deque.len();
            for _ in 0..level_length {
                let node = deque.pop_front().unwrap();
                if let Some(node) = node {
                    if node.borrow().left.is_some() {
                        deque.push_back(node.borrow().left.clone());
                    }
                    if node.borrow().right.is_some() {
                        deque.push_back(node.borrow().right.clone());
                    }
                }
            }
            depth += 1;
        }

        depth
    }

    pub fn max_depth2(root: Option<Rc<RefCell<TreeNode>>>) -> i32 {
        match root {
            Some(node) => {
                let left = TreeNode::max_depth2(node.borrow().left.clone());
                let right = TreeNode::max_depth2(node.borrow().right.clone());
                left.max(right) + 1
            }
            None => 0,
        }
    }

    pub fn min_depth(root: Option<Rc<RefCell<TreeNode>>>) -> i32 {
        if root.is_none() {
            return 0;
        }

        let mut deque: VecDeque<Option<Rc<RefCell<TreeNode>>>> = VecDeque::new();
        deque.push_back(root);

        let mut depth = 0;
        while !deque.is_empty() {
            let level_length = deque.len();
            for _ in 0..level_length {
                let node = deque.pop_front().unwrap();
                if let Some(node) = node {
                    if node.borrow().left.is_none() && node.borrow().right.is_none() {
                        return depth + 1;
                    }
                    if node.borrow().left.is_some() {
                        deque.push_back(node.borrow().left.clone());
                    }
                    if node.borrow().right.is_some() {
                        deque.push_back(node.borrow().right.clone());
                    }
                }
            }
            depth += 1;
        }

        depth
    }

    pub fn search_bst(
        root: Option<Rc<RefCell<TreeNode>>>,
        val: i32,
    ) -> Option<Rc<RefCell<TreeNode>>> {
        if root.is_none() {
            return None;
        }

        let mut r = root.clone();
        while let Some(node) = r {
            if node.borrow().val == val {
                return Some(node);
            } else if node.borrow().val > val {
                r = node.borrow().left.clone();
            } else {
                r = node.borrow().right.clone();
            }
        }

        None
    }
}

// pub fn insert_into_bst(root: Option<Rc<RefCell<TreeNode>>>,val: i32) -> Option<Rc<RefCell<TreeNode>>> {
//     if let Some(r) = &root{{
//         let mut root = r.borrow_mut();
//         if val < root.val {
//             root.left=Self::insert_into_bst(root.left.take(),val)
//         } else {
//             root.right=Self::insert_into_bst(root.right.take(),val)
//         }}
//         root
//     }else{Some(Rc::new(RefCell::new(TreeNode {left:None,right:None,val: val})))}
// }

pub fn my_pow(x: f64, n: i32) -> f64 {
    if n == 0 {
        return 1.0;
    }

    if n < 0 {
        return 1.0 / my_pow(x, -n);
    }

    if n % 2 == 0 {
        let y = my_pow(x, n / 2);
        return y * y;
    } else {
        let y = my_pow(x, n / 2);
        return y * y * x;
    }
}

pub fn climb_stairs(n: i32) -> i32 {
    if n <= 2 {
        return n;
    }

    let mut dp = vec![0; n as usize + 1];
    dp[0] = 1;
    dp[1] = 1;

    for i in 2..=n as usize {
        dp[i] = dp[i - 1] + dp[i - 2];
    }

    dp[n as usize]
}

pub fn recursion(vec: &mut Vec<String>, left: i32, right: i32, n: i32, s: String) {
    if left == n && right == n {
        vec.push(s.clone());
    }
    if left < n {
        recursion(vec, left + 1, right, n, s.clone() + "(");
    }
    if right < left {
        recursion(vec, left, right + 1, n, s.clone() + ")");
    }
}

pub fn generate_parenthesis(n: i32) -> Vec<String> {
    let mut vec: Vec<String> = Vec::new();
    recursion(&mut vec, 0, 0, n, String::from(""));

    vec
}

pub fn subsets(nums: Vec<i32>) -> Vec<Vec<i32>> {
    if nums.len() == 0 {
        return Vec::new();
    }

    let mut vecs: Vec<Vec<i32>> = Vec::new();
    let mut vec: Vec<i32> = Vec::new();
    backtrack(&mut vecs, &mut vec, &nums, 0);

    vecs
}

fn backtrack(vecs: &mut Vec<Vec<i32>>, vec: &mut Vec<i32>, nums: &Vec<i32>, index: usize) {
    vecs.push(vec.clone());

    for i in index..nums.len() {
        vec.push(nums[i]);
        backtrack(vecs, vec, nums, i + 1);
        // vec.pop();
        vec.remove(vec.len() - 1);
    }
}

pub fn combine(n: i32, k: i32) -> Vec<Vec<i32>> {
    if n < k {
        return Vec::new();
    }

    let mut vecs: Vec<Vec<i32>> = Vec::new();
    let mut vec: Vec<i32> = Vec::new();
    backtrack_combine(&mut vecs, &mut vec, n, k, 1);

    vecs
}

fn backtrack_combine(vecs: &mut Vec<Vec<i32>>, vec: &mut Vec<i32>, n: i32, k: i32, index: usize) {
    if vec.len() == k as usize {
        vecs.push(vec.clone());
        return;
    }

    let mut i = index;
    while i <= (n - (k - vec.len() as i32) + 1) as usize {
        vec.push(i as i32);
        backtrack_combine(vecs, vec, n, k, i + 1);
        vec.pop();

        i += 1;
    }
}

pub fn solve_n_queens(n: i32) -> Vec<Vec<String>> {
    if n == 0 {
        return Vec::new();
    }

    let mut board = vec![vec!['.'; n as usize]; n as usize];
    let mut res: Vec<Vec<String>> = vec![];
    backtrack_n_queens(&mut board, &mut res, n, 0);

    res
}

pub fn backtrack_n_queens(
    board: &mut Vec<Vec<char>>,
    res: &mut Vec<Vec<String>>,
    n: i32,
    row: i32,
) {
    for column in 0..n {
        if !collision(&board, n, row, column) {
            board[row as usize][column as usize] = 'Q';
            if row == n - 1 {
                res.push(board.iter().map(|vec| vec.iter().collect()).collect());
            } else {
                backtrack_n_queens(board, res, n, row + 1);
            }
            board[row as usize][column as usize] = '.';
        }
    }
}

pub fn collision(board: &Vec<Vec<char>>, n: i32, row: i32, column: i32) -> bool {
    let mut up_row = row - 1;
    let mut left_column = column - 1;
    let mut right_column = column + 1;

    while up_row >= 0 {
        if board[up_row as usize][column as usize] == 'Q' {
            return true;
        }
        if left_column >= 0 && board[up_row as usize][left_column as usize] == 'Q' {
            return true;
        }
        if right_column < n && board[up_row as usize][right_column as usize] == 'Q' {
            return true;
        }
        up_row -= 1;
        left_column -= 1;
        right_column += 1;
    }

    false
}

fn middle_search(nums: Vec<i32>, target: i32) -> i32 {
    let mut left = 0;
    let mut right = nums.len() - 1;

    while left <= right {
        let mid = (left + right) / 2;
        if nums[mid] == target {
            return mid as i32;
        } else if nums[mid] > target {
            right = mid - 1;
        } else {
            left = mid + 1;
        }
    }

    -1
}

pub fn rotate_search(nums: Vec<i32>, target: i32) -> i32 {
    if nums.len() == 0 {
        return -1;
    }

    let mut left = 0;
    let mut right = nums.len() - 1;

    while left <= right {
        let mid = (left + right) / 2;
        if nums[mid] == target {
            return mid as i32;
        } else if nums[left] <= nums[mid] {
            if target >= nums[left] && target < nums[mid] {
                right = mid - 1;
            } else {
                left = mid + 1;
            }
        } else {
            if target > nums[mid] && target <= nums[right] {
                left = mid + 1;
            } else {
                right = mid - 1;
            }
        }
    }

    -1
}

pub fn is_perfect_square(num: i32) -> bool {
    if num == 1 {
        return true;
    }

    let mut left = 1;
    let mut right = num;

    while left <= right {
        let mid = (left + right) / 2;
        let square = mid * mid;
        if square == num {
            return true;
        } else if square > num {
            right = mid - 1;
        } else {
            left = mid + 1;
        }
    }

    false
}

pub fn bubble_sort(nums: &mut Vec<i32>) {
    let mut swapped = true;
    while swapped {
        swapped = false;
        for i in 0..nums.len() - 1 {
            if nums[i] > nums[i + 1] {
                nums.swap(i, i + 1);
                swapped = true;
            }
        }
    }
}

fn insertion_sort(nums: &mut Vec<i32>) {
    for i in 1..nums.len() {
        let mut j = i;
        while j > 0 && nums[j - 1] > nums[j] {
            nums.swap(j - 1, j);
            j -= 1;
        }
    }
}

fn selection_sort(nums: &mut Vec<i32>) {
    for i in 0..nums.len() {
        let mut min_index = i;
        for j in i..nums.len() {
            if nums[j] < nums[min_index] {
                min_index = j;
            }
        }
        nums.swap(i, min_index);
    }
}

fn heap_sort(nums: &mut Vec<i32>) {
    build_heap(nums);

    for i in (0..nums.len()).rev() {
        nums.swap(0, i);
        heapify(nums, 0, i);
    }
}

fn build_heap(nums: &mut Vec<i32>) {
    let len = nums.len();
    for i in (0..len / 2).rev() {
        heapify(nums, i, len);
    }
}

fn heapify(nums: &mut Vec<i32>, i: usize, heap_size: usize) {
    let left = 2 * i + 1;
    let right = 2 * i + 2;
    let mut largest = i;

    if left < heap_size && nums[left] > nums[largest] {
        largest = left;
    }

    if right < heap_size && nums[right] > nums[largest] {
        largest = right;
    }

    if largest != i {
        nums.swap(i, largest);
        heapify(nums, largest, heap_size);
    }
}

fn merge_sort(nums: &mut Vec<i32>) {
    let len = nums.len();
    if len <= 1 {
        return;
    }

    let mid = len / 2;
    let mut left = Vec::new();
    let mut right = Vec::new();

    for i in 0..mid {
        left.push(nums[i]);
    }

    for i in mid..len {
        right.push(nums[i]);
    }

    merge_sort(&mut left);
    merge_sort(&mut right);

    merge(nums, &left, &right);
}

fn merge(nums: &mut Vec<i32>, left: &Vec<i32>, right: &Vec<i32>) {
    let mut i = 0;
    let mut j = 0;
    let mut k = 0;

    while i < left.len() && j < right.len() {
        if left[i] <= right[j] {
            nums[k] = left[i];
            i += 1;
        } else {
            nums[k] = right[j];
            j += 1;
        }
        k += 1;
    }

    while i < left.len() {
        nums[k] = left[i];
        i += 1;
        k += 1;
    }

    while j < right.len() {
        nums[k] = right[j];
        j += 1;
        k += 1;
    }
}

fn quick_sort(nums: &mut Vec<i32>) {
    let len = nums.len();
    if len <= 1 {
        return;
    }

    let pivot = nums[0];
    let mut left = Vec::new();
    let mut right = Vec::new();

    for i in 1..len {
        if nums[i] < pivot {
            left.push(nums[i]);
        } else {
            right.push(nums[i]);
        }
    }

    quick_sort(&mut left);
    quick_sort(&mut right);

    left.push(pivot);
    left.append(&mut right);

    for i in 0..len {
        nums[i] = left[i];
    }
}

pub fn coin_change(coins: Vec<i32>, amount: i32) -> i32 {
    let mut dp = vec![amount + 1; amount as usize + 1];
    dp[0] = 0;
    for coin in coins {
        (coin as usize..=amount as usize).for_each(|i| {
            dp[i] = dp[i].min(dp[i - coin as usize] + 1);
        });
    }
    if dp[amount as usize] == amount + 1 {
        -1
    } else {
        dp[amount as usize]
    }
}

pub fn length_of_lis(nums: Vec<i32>) -> i32 {
    let mut dp = vec![1; nums.len()];
    let mut max = 0;
    for i in 1..nums.len() {
        for j in 0..i {
            if nums[i] > nums[j] {
                dp[i] = dp[i].max(dp[j] + 1);
            }
        }
        max = max.max(dp[i]);
    }
    max
}

pub fn intersection(nums: Vec<Vec<i32>>) -> Vec<i32> {
    let mut result = vec![];
    let mut map = HashMap::new();
    for num in &nums {
        for i in num {
            map.insert(i, map.get(&i).unwrap_or(&0) + 1);
        }
    }

    for (k, v) in map.iter() {
        if *v == nums.len() {
            result.push(**k);
        }
    }

    result.sort();
    result
}

fn is_unique(s: String) -> bool {
    let mut map = HashMap::new();
    for c in s.chars() {
        if let Some(v) = map.get(&c) {
            if *v == 1 {
                return false;
            }
        } else {
            map.insert(c, 1);
        }
    }
    true
}

fn check_permutation(s1: String, s2: String) -> bool {
    if s1.len() != s2.len() {
        return false;
    }

    let mut map = HashMap::new();
    for c in s1.chars() {
        if let Some(v) = map.get(&c) {
            map.insert(c, *v + 1);
        } else {
            map.insert(c, 1);
        }
    }

    for c in s2.chars() {
        if let Some(v) = map.get(&c) {
            if *v == 1 {
                map.remove(&c);
            } else {
                map.insert(c, *v - 1);
            }
        } else {
            return false;
        }
    }

    map.is_empty()
}

fn check_permutation2(s1: String, s2: String) -> bool {
    let (mut a, mut b) = ([0; 26], [0; 26]);
    s1.chars().for_each(|c| a[c as usize - 'a' as usize] += 1);
    s2.chars().for_each(|c| b[c as usize - 'a' as usize] += 1);

    a == b
}

fn replace_spaces(s: String, length: i32) -> String {
    s[..length as usize].replace(" ", "%20")
}

fn can_permute_palindrome(s: String) -> bool {
    let mut map = HashMap::new();
    for c in s.chars() {
        if let Some(v) = map.get(&c) {
            if *v == 1 {
                map.remove(&c);
            } else {
                map.insert(c, *v + 1);
            }
        } else {
            map.insert(c, 1);
        }
    }

    map.values().filter(|&v| *v % 2 == 1).count() <= 1
}

fn remove_duplicate_nodes(head: Option<Box<ListNode>>) -> Option<Box<ListNode>> {
    let mut cur = head;
    let mut set = HashSet::new();
    let mut v = Vec::new();
    while let Some(node) = cur {
        if !set.contains(&node.val) {
            set.insert(node.val);
            v.push(node.val);
        }
        cur = node.next;
    }

    v.reverse();
    let mut prev = None;
    for i in v {
        let node = Box::new(ListNode { val: i, next: prev });
        prev = Some(node);
    }

    prev
}

pub fn is_fliped_string(s1: String, s2: String) -> bool {
    s1.len() == s2.len() && s1.repeat(2).contains(&s2)
}

pub fn set_zeroes(matrix: &mut Vec<Vec<i32>>) {
    let mut row = vec![false; matrix.len()];
    let mut col = vec![false; matrix[0].len()];

    for i in 0..matrix.len() {
        for j in 0..matrix[0].len() {
            if matrix[i][j] == 0 {
                row[i] = true;
                col[j] = true;
            }
        }
    }

    for i in 0..matrix.len() {
        for j in 0..matrix[0].len() {
            if row[i] || col[j] {
                matrix[i][j] = 0;
            }
        }
    }
}

pub fn rotate(matrix: &mut Vec<Vec<i32>>) {
    if matrix.is_empty() {
        return;
    }
    let col = matrix.len();
    let row = matrix[0].len();

    for i in 0..row / 2 {
        for j in 0..col {
            let tmp = matrix[i][j];
            matrix[i][j] = matrix[row - 1 - i][j];
            matrix[row - 1 - i][j] = tmp;
        }
    }

    for i in 0..row {
        for j in 0..i {
            let tmp = matrix[i][j];
            matrix[i][j] = matrix[j][i];
            matrix[j][i] = tmp;
        }
    }
}

pub fn compress_string(s: String) -> String {
    if s.is_empty() {
        return s;
    }
    let mut res = String::new();
    let mut count = 1;
    let mut prev = s.chars().next().unwrap();
    for c in s.chars().skip(1) {
        if c == prev {
            count += 1;
        } else {
            res.push(prev);
            res.push_str(&count.to_string());
            count = 1;
            prev = c;
        }
    }
    res.push(prev);
    res.push_str(&count.to_string());
    if res.len() >= s.len() {
        s
    } else {
        res
    }
}

fn one_edit_away(s1: String, s2: String) -> bool {
    let mut i = 0;
    let mut j = 0;
    let mut count = 0;
    while i < s1.len() && j < s2.len() {
        if s1.chars().nth(i) == s2.chars().nth(j) {
            i += 1;
            j += 1;
        } else {
            count += 1;
            if count > 1 {
                return false;
            }
            if s1.len() > s2.len() {
                i += 1;
            } else if s1.len() < s2.len() {
                j += 1;
            } else {
                i += 1;
                j += 1;
            }
        }
    }
    count + (s1.len() - i) + (s2.len() - j) <= 2
}

fn kth_to_last(head: Option<Box<ListNode>>, k: i32) -> i32 {
    let mut cur = &head;
    let mut count = 0;
    while let Some(node) = cur {
        count += 1;
        cur = &node.next;
    }
    if count < k {
        return -1;
    }
    let mut cur = head;
    for _ in 0..count - k {
        cur = cur.unwrap().next;
    }

    cur.unwrap().val
}

fn delete_middle_node(node: &mut Box<ListNode>) {
    let next = node.next.take();
    *node = next.unwrap();
}

fn addTwoNumbers(l1: Option<Box<ListNode>>, l2: Option<Box<ListNode>>) -> Option<Box<ListNode>> {
    let mut l1 = &l1;
    let mut l2 = &l2;
    let mut head = None;
    let mut cur = &mut head;
    let mut carry = 0;
    while l1.is_some() || l2.is_some() || carry != 0 {
        let mut sum = carry;
        if let Some(n1) = l1.as_ref() {
            sum += n1.val;
            l1 = &n1.next;
        }
        if let Some(n2) = l2.as_ref() {
            sum += n2.val;
            l2 = &n2.next;
        }
        carry = sum / 10;
        *cur = Option::Some(Box::new(ListNode::new(sum % 10)));
        cur = &mut cur.as_mut().unwrap().next;
    }

    head
}

pub fn partition(head: Option<Box<ListNode>>, x: i32) -> Option<Box<ListNode>> {
    let mut lower = Some(Box::new(ListNode::new(0)));
    let mut higher = Some(Box::new(ListNode::new(0)));
    let mut lower_cur = lower.as_mut();
    let mut higher_cur = higher.as_mut();
    let mut cur = head;

    while let Some(mut node) = cur {
        let mut next = node.next.take();
        if node.val < x {
            lower_cur.as_mut().unwrap().next = Some(node);
            lower_cur = lower_cur.unwrap().next.as_mut();
        } else {
            higher_cur.as_mut().unwrap().next = Some(node);
            higher_cur = higher_cur.unwrap().next.as_mut();
        }
        cur = next;
    }

    lower_cur.as_mut().unwrap().next = higher.unwrap().next;

    lower.unwrap().next.take()
}

pub fn is_palindrome(head: Option<Box<ListNode>>) -> bool {
    let mut cur = head;
    let mut vals = Vec::new();
    while let Some(node) = cur {
        vals.push(node.val);
        cur = node.next;
    }

    let mut i = 0;
    let mut j = vals.len() - 1;
    while i < j {
        if vals[i] != vals[j] {
            return false;
        }
        i += 1;
        j -= 1;
    }

    true
}

fn get_intersection_node(
    l1: Option<Box<ListNode>>,
    l2: Option<Box<ListNode>>,
) -> Option<Box<ListNode>> {
    let mut l1 = &l1;
    let mut l2 = &l2;
    // let mut dummy = Some(Box::new(ListNode::new(0)));
    while let Some(n1) = l1 {
        let mut dummy = l2;
        while let Some(n2) = dummy {
            if n1 == n2 {
                return Some(Box::new(ListNode::new(n1.val)));
            }
            dummy = &n2.next;
        }
        l1 = &n1.next;
    }

    None
}

struct TripleInOne {
    d: Vec<i32>,
    i: [usize; 3],
}

impl TripleInOne {
    fn new(n: i32) -> Self {
        TripleInOne {
            d: vec![0; n as usize],
            i: [0, 1, 2],
        }
    }

    fn push(&mut self, stack_num: i32, value: i32) {
        let n = stack_num as usize;
        if self.i[n] < self.d.len() {
            self.d[self.i[n]] = value;
            self.i[n] += 3;
        }
    }

    fn pop(&mut self, stack_num: i32) -> i32 {
        match stack_num as usize {
            n if self.i[n] >= 3 => {
                self.i[n] -= 3;
                self.d[self.i[n]]
            }
            _ => -1,
        }
    }
}

struct StackOfPlates {
    inner: Vec<Vec<i32>>,
    cap: usize,
}

impl StackOfPlates {
    fn new(capacity: i32) -> Self {
        StackOfPlates {
            inner: Vec::new(),
            cap: capacity as usize,
        }
    }

    fn push(&mut self, val: i32) {
        if self.cap == 0 {
            return;
        }
        let last = self.inner.last_mut();
        if let Some(v) = last {
            if v.len() < self.cap {
                v.push(val);
                return;
            }
        }
        let mut stack = Vec::with_capacity(self.cap);
        stack.push(val);
        self.inner.push(stack);
    }

    fn pop(&mut self) -> i32 {
        let last = self.inner.last_mut();
        if let Some(v) = last {
            if let Some(val) = v.pop() {
                if v.is_empty() {
                    self.inner.pop();
                }
                return val;
            }
        }
        -1
    }

    fn pop_at(&mut self, index: i32) -> i32 {
        let stack = self.inner.get_mut(index as usize);
        if let Some(vec_ref) = stack {
            if let Some(val) = vec_ref.pop() {
                if vec_ref.is_empty() {
                    self.inner.remove(index as usize);
                }
                return val;
            }
        }
        -1
    }
}

struct MyQueue {
    vec: VecDeque<i32>,
}

impl MyQueue {
    /** Initialize your data structure here. */
    fn new() -> Self {
        MyQueue {
            vec: VecDeque::new(),
        }
    }

    /** Push element x to the back of queue. */
    fn push(&mut self, x: i32) {
        self.vec.push_back(x);
    }

    /** Removes the element from in front of queue and returns that element. */
    fn pop(&mut self) -> i32 {
        if self.vec.len() > 0 {
            self.vec.pop_front().unwrap()
        } else {
            -1
        }
    }

    /** Get the front element. */
    fn peek(&self) -> i32 {
        let l = self.vec.len();
        if self.vec.len() > 0 {
            self.vec[0]
        } else {
            -1
        }
    }

    /** Returns whether the queue is empty. */
    fn empty(&self) -> bool {
        self.vec.is_empty()
    }
}

pub fn find_whether_exists_path(n: i32, graph: Vec<Vec<i32>>, start: i32, target: i32) -> bool {
    let mut visited = vec![false; n as usize];
    visited[start as usize] = true;
    let mut queue = VecDeque::new();
    queue.push_back(start);
    while !queue.is_empty() {
        let cur = queue.pop_front().unwrap();
        if cur == target {
            return true;
        }
        for next in graph[cur as usize].iter() {
            if !visited[*next as usize] {
                visited[*next as usize] = true;
                queue.push_back(*next);
            }
        }
    }

    false
}

pub fn is_balanced(root: Option<Rc<RefCell<TreeNode>>>) -> bool {
    fn post_order(node: Option<Rc<RefCell<TreeNode>>>) -> i32 {
        match node {
            None => 0,
            Some(node) => {
                let left = post_order(node.borrow().left.clone());
                let right = post_order(node.borrow().right.clone());
                if left == -1 || right == -1 || (left - right).abs() > 1 {
                    -1
                } else {
                    1 + left + right
                }
            }
        }
    }

    post_order(root) != -1
}

fn sorted_array_to_bst(nums: Vec<i32>) -> Option<Rc<RefCell<TreeNode>>> {
    if nums.is_empty() {
        return None;
    }
    let mid = nums.len() / 2;
    let mut root = TreeNode::new(nums[mid]);
    root.left = sorted_array_to_bst(nums[..mid].to_vec());
    root.right = sorted_array_to_bst(nums[mid + 1..].to_vec());

    Some(Rc::new(RefCell::new(root)))
}

struct RecentCounter {
    q: VecDeque<i32>,
}

impl RecentCounter {
    fn new() -> Self {
        RecentCounter { q: VecDeque::new() }
    }

    fn ping(&mut self, t: i32) -> i32 {
        self.q.push_back(t);
        while self.q.len() > 0 && t - self.q.front().unwrap() > 3000 {
            self.q.pop_front();
        }

        self.q.len() as i32
    }
}

pub fn min_mutation(start: &str, end: &str, bank: Vec<&str>) -> i32 {
    let mut bank = bank;
    bank.sort();
    let mut queue = VecDeque::new();
    queue.push_back(start);
    let mut visited = vec![false; bank.len()];
    let mut step = 0;
    while !queue.is_empty() {
        let cur = queue.pop_front().unwrap();
        if cur == end {
            return step;
        }
        for i in 0..bank.len() {
            if !visited[i] && is_one_difference(&cur, &bank[i]) {
                visited[i] = true;
                queue.push_back(bank[i].clone());
            }
        }
        step += 1;
    }

    -1
}

fn is_one_difference(s: &str, t: &str) -> bool {
    let mut diff = 0;
    for i in 0..s.len() {
        if s.as_bytes()[i] != t.as_bytes()[i] {
            diff += 1;
        }
    }

    diff == 1
}

#[cfg(test)]
mod tests {
    use super::*;
    use std::borrow::Borrow;

    #[test]
    fn test_move_zeroes() {
        let mut nums = vec![0, 1, 0, 3, 12];
        move_zeroes(&mut nums);
        assert_eq!(nums, vec![1, 3, 12, 0, 0]);
    }

    #[test]
    fn test_plus_one() {
        assert_eq!(plus_one(vec![1, 2, 3]), vec![1, 2, 4]);
        assert_eq!(plus_one(vec![4, 3, 2, 1]), vec![4, 3, 2, 2]);
        assert_eq!(plus_one(vec![9, 9, 9, 9]), vec![1, 0, 0, 0, 0]);
    }

    #[test]
    fn test_remove_duplicates() {
        let mut nums = vec![0, 0, 1, 1, 1, 2, 2, 3, 3, 4];
        assert_eq!(remove_duplicates(&mut nums), 5);
        // assert_eq!(nums, vec![0, 1, 2, 3, 4]);
    }

    #[test]
    fn test_min_stack() {
        let mut min_stack = MinStack::new();
        min_stack.push(2);
        min_stack.push(0);
        min_stack.push(3);
        min_stack.push(0);
        assert_eq!(min_stack.get_min(), 0);
        min_stack.pop();
        assert_eq!(min_stack.top(), 3);
        assert_eq!(min_stack.get_min(), 0);
    }

    #[test]
    fn test_is_valid() {
        assert_eq!(is_valid("()".to_string()), true);
        assert_eq!(is_valid("()[]{}".to_string()), true);
        assert_eq!(is_valid("(]".to_string()), false);
        assert_eq!(is_valid("([)]".to_string()), false);
        assert_eq!(is_valid("{[]}".to_string()), true);
    }

    #[test]
    fn test_max_sliding_window() {
        assert_eq!(
            max_sliding_window(vec![1, 3, -1, -3, 5, 3, 6, 7], 3),
            vec![3, 3, 5, 5, 6, 7]
        );
    }

    #[test]
    fn test_two_sum() {
        assert_eq!(two_sum(vec![2, 7, 11, 15], 9), vec![0, 1]);
    }

    #[test]
    fn test_is_anagram() {
        assert_eq!(
            is_anagram("anagram".to_string(), "nagaram".to_string()),
            true
        );
        assert_eq!(is_anagram("rat".to_string(), "car".to_string()), false);
    }

    // #[test]
    fn test_group_anagram() {
        assert_eq!(
            group_anagram(vec![
                "eat".to_string(),
                "tea".to_string(),
                "tan".to_string(),
                "ate".to_string(),
                "nat".to_string(),
                "bat".to_string()
            ]),
            vec![
                vec!["tan".to_string(), "nat".to_string()],
                vec!["bat".to_string()],
                vec!["ate".to_string(), "eat".to_string(), "tea".to_string()],
            ]
        );
    }

    // #[test]
    // fn test_reverse_list() {
    //     let mut head = Some(Box::new(ListNode {
    //         val: 1,
    //         next: Some(Box::new(ListNode {
    //             val: 2,
    //             next: Some(Box::new(ListNode {
    //                 val: 3,
    //                 next: Some(Box::new(ListNode { val: 4, next: None })),
    //             })),
    //         })),
    //     }));
    //     let mut new_head = reverse_list(head);
    //     assert_eq!(new_head.as_ref().unwrap().val, 4);
    //     assert_eq!(new_head.unwrap().next.unwrap().val, 3);
    // }
    //
    // #[test]
    // fn test_middle_node() {
    //     let mut head = Some(Box::new(ListNode {
    //         val: 1,
    //         next: Some(Box::new(ListNode {
    //             val: 2,
    //             next: Some(Box::new(ListNode {
    //                 val: 3,
    //                 next: Some(Box::new(ListNode { val: 4, next: None })),
    //             })),
    //         })),
    //     }));
    //     assert_eq!(middle_node(head).unwrap().val, 3);
    // }

    #[test]
    fn test_merge_two_lists() {
        let mut l1 = Some(Box::new(ListNode {
            val: 1,
            next: Some(Box::new(ListNode {
                val: 2,
                next: Some(Box::new(ListNode { val: 4, next: None })),
            })),
        }));
        let mut l2 = Some(Box::new(ListNode {
            val: 1,
            next: Some(Box::new(ListNode {
                val: 3,
                next: Some(Box::new(ListNode { val: 4, next: None })),
            })),
        }));
        let new_head = merge_two_lists(l1, l2);
        assert_eq!(new_head.as_ref().unwrap().val, 1);
        assert_eq!(new_head.as_ref().unwrap().next.as_ref().unwrap().val, 1);
        assert_eq!(
            new_head
                .as_ref()
                .unwrap()
                .next
                .as_ref()
                .unwrap()
                .next
                .as_ref()
                .unwrap()
                .val,
            2
        );
        // assert_eq!(new_head.unwrap().next.unwrap().next.unwrap().next.unwrap().val, 3);
        // assert_eq!(new_head.unwrap().next.unwrap().next.unwrap().next.unwrap().next.unwrap().val, 4);
        // assert_eq!(new_head.unwrap().next.unwrap().next.unwrap().next.unwrap().next.unwrap().next.unwrap().val, 4);
    }

    #[test]
    fn test_remove_nth_from_end() {
        let mut head = Some(Box::new(ListNode {
            val: 1,
            next: Some(Box::new(ListNode {
                val: 2,
                next: Some(Box::new(ListNode {
                    val: 3,
                    next: Some(Box::new(ListNode {
                        val: 4,
                        next: Some(Box::new(ListNode { val: 5, next: None })),
                    })),
                })),
            })),
        }));
        let new_head = remove_nth_from_end(head, 2);
        assert_eq!(new_head.as_ref().unwrap().val, 1);
        // assert_eq!(new_head.as_ref().unwrap().next.as_ref().unwrap().val, 2);
        // assert_eq!(new_head.as_ref().unwrap().next.as_ref().unwrap().next.as_ref().unwrap().val, 4);
        // assert_eq!(new_head.as_ref().unwrap().next.as_ref().unwrap().next.as_ref().unwrap().next.as_ref().unwrap().val, 5);
    }

    #[test]
    fn test_preorder_traversal() {
        let mut root = Option::from(Rc::new(RefCell::new(TreeNode {
            val: 1,
            left: Some(Rc::new(RefCell::new(TreeNode {
                val: 2,
                left: None,
                right: None,
            }))),
            right: Some(Rc::new(RefCell::new(TreeNode {
                val: 3,
                left: None,
                right: None,
            }))),
        })));
        let result = TreeNode::preorder_traversal(root);
        assert_eq!(result, vec![1, 2, 3]);
    }

    #[test]
    fn test_level_order() {
        let mut root = Option::from(Rc::new(RefCell::new(TreeNode {
            val: 1,
            left: Some(Rc::new(RefCell::new(TreeNode {
                val: 2,
                left: None,
                right: None,
            }))),
            right: Some(Rc::new(RefCell::new(TreeNode {
                val: 3,
                left: None,
                right: None,
            }))),
        })));
        let result = TreeNode::level_order(root);
        assert_eq!(result, vec![vec![1], vec![2, 3]]);
    }

    #[test]
    fn test_my_pow() {
        assert_eq!(my_pow(2.0, 10), 1024.0);
        assert_eq!(my_pow(2.0, -2), 0.25);
    }

    #[test]
    fn test_climb_stairs() {
        assert_eq!(climb_stairs(2), 2);
        assert_eq!(climb_stairs(3), 3);
        assert_eq!(climb_stairs(4), 5);
    }

    #[test]
    fn test_generate_parenthesis() {
        let result = generate_parenthesis(3);
        assert_eq!(
            result,
            vec!["((()))", "(()())", "(())()", "()(())", "()()()"]
        );
    }

    #[test]
    fn test_subsets() {
        let result = subsets(vec![1, 2, 3]);
        assert_eq!(
            result,
            vec![
                vec![],
                vec![1],
                vec![1, 2],
                vec![1, 2, 3],
                vec![1, 3],
                vec![2],
                vec![2, 3],
                vec![3]
            ]
        );
    }

    #[test]
    fn test_combine() {
        let result = combine(4, 2);
        assert_eq!(
            result,
            vec![
                vec![1, 2],
                vec![1, 3],
                vec![1, 4],
                vec![2, 3],
                vec![2, 4],
                vec![3, 4]
            ]
        );
    }

    #[test]
    fn test_slove_n_queens() {
        let result = solve_n_queens(4);
        assert_eq!(
            result,
            vec![
                vec![".Q..", "...Q", "Q...", "..Q."],
                vec!["..Q.", "Q...", "...Q", ".Q.."]
            ]
        );
    }

    #[test]
    fn test_middle_search() {
        let result = middle_search(vec![1, 2, 3, 4, 5, 6, 7, 8, 9, 10], 6);
        assert_eq!(result, 5);
    }

    #[test]
    fn test_rotate_search() {
        let result = rotate_search(vec![1, 2, 3, 4, 5, 6, 7, 8, 9, 10], 6);
        assert_eq!(result, 5);
    }

    #[test]
    fn test_is_perfect_square() {
        assert_eq!(is_perfect_square(1), true);
        assert_eq!(is_perfect_square(4), true);
        assert_eq!(is_perfect_square(16), true);
        assert_eq!(is_perfect_square(14), false);
    }

    #[test]
    fn test_max_depth() {
        let mut root = Option::from(Rc::new(RefCell::new(TreeNode {
            val: 1,
            left: Some(Rc::new(RefCell::new(TreeNode {
                val: 2,
                left: None,
                right: None,
            }))),
            right: Some(Rc::new(RefCell::new(TreeNode {
                val: 3,
                left: None,
                right: None,
            }))),
        })));
        assert_eq!(TreeNode::max_depth(root), 2);
    }

    #[test]
    fn test_max_depth2() {
        let mut root = Option::from(Rc::new(RefCell::new(TreeNode {
            val: 1,
            left: Some(Rc::new(RefCell::new(TreeNode {
                val: 2,
                left: None,
                right: None,
            }))),
            right: Some(Rc::new(RefCell::new(TreeNode {
                val: 3,
                left: None,
                right: Some(Rc::new(RefCell::new(TreeNode {
                    val: 4,
                    left: None,
                    right: None,
                }))),
            }))),
        })));
        assert_eq!(TreeNode::max_depth2(root), 3);
    }

    #[test]
    fn test_min_depth() {
        let mut root = Option::from(Rc::new(RefCell::new(TreeNode {
            val: 1,
            left: Some(Rc::new(RefCell::new(TreeNode {
                val: 2,
                left: None,
                right: None,
            }))),
            right: Some(Rc::new(RefCell::new(TreeNode {
                val: 3,
                left: None,
                right: Some(Rc::new(RefCell::new(TreeNode {
                    val: 4,
                    left: None,
                    right: None,
                }))),
            }))),
        })));
        assert_eq!(TreeNode::min_depth(root), 2);
    }

    #[test]
    fn test_search_bst() {
        let mut root = Option::from(Rc::new(RefCell::new(TreeNode {
            val: 4,
            left: Some(Rc::new(RefCell::new(TreeNode {
                val: 2,
                left: Some(Rc::new(RefCell::new(TreeNode {
                    val: 1,
                    left: None,
                    right: None,
                }))),
                right: Some(Rc::new(RefCell::new(TreeNode {
                    val: 3,
                    left: None,
                    right: None,
                }))),
            }))),
            right: Some(Rc::new(RefCell::new(TreeNode {
                val: 6,
                left: Some(Rc::new(RefCell::new(TreeNode {
                    val: 5,
                    left: None,
                    right: None,
                }))),
                right: Some(Rc::new(RefCell::new(TreeNode {
                    val: 7,
                    left: None,
                    right: None,
                }))),
            }))),
        })));
        assert_eq!(
            TreeNode::search_bst(root, 5),
            Some(Rc::new(RefCell::new(TreeNode {
                val: 5,
                left: None,
                right: None,
            })))
        );
    }

    #[test]
    fn test_bubble_sort() {
        let mut arr = vec![5, 4, 3, 2, 1];
        bubble_sort(&mut arr);
        assert_eq!(arr, vec![1, 2, 3, 4, 5]);
    }

    #[test]
    fn test_insertion_sort() {
        let mut arr = vec![5, 4, 3, 2, 1];
        insertion_sort(&mut arr);
        assert_eq!(arr, vec![1, 2, 3, 4, 5]);
    }

    #[test]
    fn test_selection_sort() {
        let mut arr = vec![5, 4, 3, 2, 1];
        selection_sort(&mut arr);
        assert_eq!(arr, vec![1, 2, 3, 4, 5]);
    }

    #[test]
    fn test_heap_sort() {
        let mut arr = vec![5, 4, 3, 2, 1];
        heap_sort(&mut arr);
        assert_eq!(arr, vec![1, 2, 3, 4, 5]);
    }

    #[test]
    fn test_merge_sort() {
        let mut arr = vec![5, 4, 3, 2, 1];
        merge_sort(&mut arr);
        assert_eq!(arr, vec![1, 2, 3, 4, 5]);
    }

    #[test]
    fn test_quick_sort() {
        let mut arr = vec![5, 4, 3, 2, 1];
        quick_sort(&mut arr);
        assert_eq!(arr, vec![1, 2, 3, 4, 5]);
    }

    #[test]
    fn test_coin_change() {
        let res = vec![1, 2, 5];
        assert_eq!(coin_change(res, 11), 3);
    }

    #[test]
    fn test_length_of_lis() {
        let res = vec![1, 3, 5, 4, 7, 9, 2, 5];
        assert_eq!(length_of_lis(res), 5);
    }

    #[test]
    fn test_intersection() {
        let res = vec![vec![3, 1, 2, 4, 5], vec![1, 2, 3, 4], vec![3, 4, 5, 6]];
        assert_eq!(intersection(res), vec![3, 4]);
    }

    #[test]
    fn test_is_unique() {
        assert_eq!(is_unique("abcdefg".to_string()), true);
        assert_eq!(is_unique("abcdefgf".to_string()), false);
    }

    #[test]
    fn test_check_permutation() {
        assert_eq!(
            check_permutation("abcdefg".to_string(), "abcdefg".to_string()),
            true
        );
        assert_eq!(
            check_permutation("abcdefg".to_string(), "abcdefgf".to_string()),
            false
        );
    }

    #[test]
    fn test_check_permutation2() {
        assert_eq!(
            check_permutation2("abcdefg".to_string(), "abcdefg".to_string()),
            true
        );
        assert_eq!(
            check_permutation2("abcdefg".to_string(), "abcdefgf".to_string()),
            false
        );
    }

    #[test]
    fn test_replace_spaces() {
        assert_eq!(
            replace_spaces("Mr John Smith    ".to_string(), 13),
            "Mr%20John%20Smith".to_string()
        );
    }

    #[test]
    fn test_can_permute_palindrome() {
        assert_eq!(can_permute_palindrome("tactcoa".to_string()), true);
        assert_eq!(can_permute_palindrome("tactcoaf".to_string()), false);
    }

    #[test]
    fn test_remove_duplicate_nodes() {
        let mut head = Option::Some(Box::new(ListNode::new(1)));
        head.as_mut().unwrap().next = Option::Some(Box::new(ListNode::new(2)));
        head.as_mut().unwrap().next.as_mut().unwrap().next =
            Option::Some(Box::new(ListNode::new(3)));
        head.as_mut()
            .unwrap()
            .next
            .as_mut()
            .unwrap()
            .next
            .as_mut()
            .unwrap()
            .next = Option::Some(Box::new(ListNode::new(3)));
        let res = remove_duplicate_nodes(head);

        let mut cur = res;
        assert_eq!(cur.as_ref().unwrap().val, 1);
        assert_eq!(cur.as_ref().unwrap().next.as_ref().unwrap().val, 2);
        assert_eq!(
            cur.as_ref()
                .unwrap()
                .next
                .as_ref()
                .unwrap()
                .next
                .as_ref()
                .unwrap()
                .val,
            3
        );
        assert_eq!(
            cur.as_ref()
                .unwrap()
                .next
                .as_ref()
                .unwrap()
                .next
                .as_ref()
                .unwrap()
                .next
                .is_none(),
            true
        );
    }

    #[test]
    fn test_is_fliped_string() {
        assert_eq!(
            is_fliped_string("waterbottle".to_string(), "erbottlewat".to_string()),
            true
        );
        assert_eq!(
            is_fliped_string("abcdefg".to_string(), "cbaefg".to_string()),
            false
        );
    }

    #[test]
    fn test_set_zeroes() {
        let mut matrix = vec![vec![1, 1, 1], vec![1, 0, 1], vec![1, 1, 1]];
        set_zeroes(&mut matrix);
        assert_eq!(matrix, vec![vec![1, 0, 1], vec![0, 0, 0], vec![1, 0, 1]]);
    }

    #[test]
    fn test_rotate() {
        let mut matrix = vec![vec![1, 2, 3], vec![4, 5, 6], vec![7, 8, 9]];
        rotate(&mut matrix);
        assert_eq!(matrix, vec![vec![7, 4, 1], vec![8, 5, 2], vec![9, 6, 3]]);
    }

    #[test]
    fn test_compress_string() {
        assert_eq!(
            compress_string("aabcccccaaa".to_string()),
            "a2b1c5a3".to_string()
        );
        assert_eq!(compress_string("a".to_string()), "a".to_string());
        assert_eq!(compress_string("".to_string()), "".to_string());
    }

    #[test]
    fn test_one_edit_away() {
        assert_eq!(one_edit_away("pale".to_string(), "ple".to_string()), true);
        assert_eq!(one_edit_away("pales".to_string(), "pale".to_string()), true);
        assert_eq!(one_edit_away("pale".to_string(), "bale".to_string()), true);
        assert_eq!(one_edit_away("pale".to_string(), "bake".to_string()), false);
    }

    #[test]
    fn test_kth_to_last() {
        let mut head = Option::Some(Box::new(ListNode::new(1)));
        head.as_mut().unwrap().next = Option::Some(Box::new(ListNode::new(2)));
        head.as_mut().unwrap().next.as_mut().unwrap().next =
            Option::Some(Box::new(ListNode::new(3)));
        head.as_mut()
            .unwrap()
            .next
            .as_mut()
            .unwrap()
            .next
            .as_mut()
            .unwrap()
            .next = Option::Some(Box::new(ListNode::new(4)));
        head.as_mut()
            .unwrap()
            .next
            .as_mut()
            .unwrap()
            .next
            .as_mut()
            .unwrap()
            .next
            .as_mut()
            .unwrap()
            .next = Option::Some(Box::new(ListNode::new(5)));
        assert_eq!(kth_to_last(head, 3), 3);
    }

    #[test]
    fn test_delete_middle_node() {
        let mut head = Option::Some(Box::new(ListNode::new(1)));
        head.as_mut().unwrap().next = Option::Some(Box::new(ListNode::new(2)));
        head.as_mut().unwrap().next.as_mut().unwrap().next =
            Option::Some(Box::new(ListNode::new(3)));
        head.as_mut()
            .unwrap()
            .next
            .as_mut()
            .unwrap()
            .next
            .as_mut()
            .unwrap()
            .next = Option::Some(Box::new(ListNode::new(4)));
        head.as_mut()
            .unwrap()
            .next
            .as_mut()
            .unwrap()
            .next
            .as_mut()
            .unwrap()
            .next
            .as_mut()
            .unwrap()
            .next = Option::Some(Box::new(ListNode::new(5)));
        delete_middle_node(head.as_mut().unwrap().next.as_mut().unwrap());
        assert_eq!(
            head,
            Option::Some(Box::new(ListNode {
                val: 1,
                next: Option::Some(Box::new(ListNode {
                    val: 3,
                    next: Option::Some(Box::new(ListNode {
                        val: 4,
                        next: Option::Some(Box::new(ListNode {
                            val: 5,
                            next: Option::None
                        }))
                    }))
                }))
            }))
        );

        let mut head = Option::Some(Box::new(ListNode::new(1)));
        head.as_mut().unwrap().next = Option::Some(Box::new(ListNode::new(2)));
        head.as_mut().unwrap().next.as_mut().unwrap().next =
            Option::Some(Box::new(ListNode::new(3)));
        head.as_mut()
            .unwrap()
            .next
            .as_mut()
            .unwrap()
            .next
            .as_mut()
            .unwrap()
            .next = Option::Some(Box::new(ListNode::new(4)));
        head.as_mut()
            .unwrap()
            .next
            .as_mut()
            .unwrap()
            .next
            .as_mut()
            .unwrap()
            .next
            .as_mut()
            .unwrap()
            .next = Option::Some(Box::new(ListNode::new(5)));
        delete_middle_node(
            head.as_mut()
                .unwrap()
                .next
                .as_mut()
                .unwrap()
                .next
                .as_mut()
                .unwrap(),
        );
        assert_eq!(
            head,
            Option::Some(Box::new(ListNode {
                val: 1,
                next: Option::Some(Box::new(ListNode {
                    val: 2,
                    next: Option::Some(Box::new(ListNode {
                        val: 4,
                        next: Option::Some(Box::new(ListNode {
                            val: 5,
                            next: Option::None
                        }))
                    }))
                }))
            }))
        );
    }

    #[test]
    fn test_addTwoNumbers() {
        let mut l1 = Option::Some(Box::new(ListNode::new(2)));
        l1.as_mut().unwrap().next = Option::Some(Box::new(ListNode::new(4)));
        l1.as_mut().unwrap().next.as_mut().unwrap().next = Option::Some(Box::new(ListNode::new(3)));
        let mut l2 = Option::Some(Box::new(ListNode::new(5)));
        l2.as_mut().unwrap().next = Option::Some(Box::new(ListNode::new(6)));
        l2.as_mut().unwrap().next.as_mut().unwrap().next = Option::Some(Box::new(ListNode::new(4)));
        let mut result = addTwoNumbers(l1, l2);
        assert_eq!(result.as_mut().unwrap().val, 7);
        assert_eq!(result.as_mut().unwrap().next.as_mut().unwrap().val, 0);
        assert_eq!(
            result
                .as_mut()
                .unwrap()
                .next
                .as_mut()
                .unwrap()
                .next
                .as_mut()
                .unwrap()
                .val,
            8
        );
        assert_eq!(
            result
                .as_mut()
                .unwrap()
                .next
                .as_mut()
                .unwrap()
                .next
                .as_mut()
                .unwrap()
                .next,
            None
        );
    }

    #[test]
    fn test_partition() {
        assert_eq!(
            partition(
                Option::Some(Box::new(ListNode {
                    val: 1,
                    next: Option::Some(Box::new(ListNode {
                        val: 4,
                        next: Option::Some(Box::new(ListNode {
                            val: 2,
                            next: Option::None
                        }))
                    }))
                })),
                3
            ),
            Option::Some(Box::new(ListNode {
                val: 1,
                next: Option::Some(Box::new(ListNode {
                    val: 2,
                    next: Option::Some(Box::new(ListNode {
                        val: 4,
                        next: Option::None
                    }))
                }))
            }))
        );
    }

    #[test]
    fn test_is_palindrome() {
        assert_eq!(
            is_palindrome(Option::Some(Box::new(ListNode {
                val: 1,
                next: Option::Some(Box::new(ListNode {
                    val: 2,
                    next: Option::Some(Box::new(ListNode {
                        val: 2,
                        next: Option::Some(Box::new(ListNode {
                            val: 1,
                            next: Option::None
                        }))
                    }))
                }))
            }))),
            true
        );
        assert_eq!(
            is_palindrome(Option::Some(Box::new(ListNode {
                val: 1,
                next: Option::Some(Box::new(ListNode {
                    val: 2,
                    next: Option::Some(Box::new(ListNode {
                        val: 2,
                        next: Option::Some(Box::new(ListNode {
                            val: 1,
                            next: Option::Some(Box::new(ListNode {
                                val: 3,
                                next: Option::None
                            }))
                        }))
                    }))
                }))
            }))),
            false
        );
    }

    #[test]
    fn test_get_intersection_node() {
        let mut l1 = Option::Some(Box::new(ListNode::new(1)));
        l1.as_mut().unwrap().next = Option::Some(Box::new(ListNode::new(2)));
        l1.as_mut().unwrap().next.as_mut().unwrap().next = Option::Some(Box::new(ListNode::new(7)));
        let mut l2 = Option::Some(Box::new(ListNode::new(4)));
        l2.as_mut().unwrap().next = Option::Some(Box::new(ListNode::new(2)));
        l2.as_mut().unwrap().next.as_mut().unwrap().next = Option::Some(Box::new(ListNode::new(6)));
        l2.as_mut()
            .unwrap()
            .next
            .as_mut()
            .unwrap()
            .next
            .as_mut()
            .unwrap()
            .next = Option::Some(Box::new(ListNode::new(7)));
        assert_eq!(
            get_intersection_node(l1, l2),
            Option::Some(Box::new(ListNode::new(7)))
        );
    }

    #[test]
    fn test_find_whether_exists_path() {
        let mut graph = vec![vec![]; 6];
        graph[0] = vec![0, 1];
        graph[1] = vec![0, 2];
        graph[2] = vec![1, 2];
        graph[3] = vec![1, 2];
        assert_eq!(find_whether_exists_path(3, graph, 0, 2), true);
    }

    #[test]
    fn test_min_mutation() {
        assert_eq!(min_mutation("AACCGGTT", "AACCGGTA", vec!["AACCGGTA"]), 1);
        assert_eq!(min_mutation("AACCGGTT", "AAACGGTA", vec!["AACCGGTA", "AACCGCTA", "AAACGGTA"]), 2);
    }
}
