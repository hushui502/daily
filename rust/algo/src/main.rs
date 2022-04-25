use std::borrow::BorrowMut;
use std::cell::RefCell;
use std::collections::{HashMap, VecDeque};
use std::fmt::format;
use std::process::id;
use std::rc::Rc;

fn main() {
    println!("Hello, world!");
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

pub fn merge_two_lists(l1: Option<Box<ListNode>>, l2: Option<Box<ListNode>>) -> Option<Box<ListNode>> {
    match (l1, l2) {
        (Some(n1), Some(n2)) =>
            if n1.val < n2.val {
                Some(Box::new(ListNode { val: n1.val, next: merge_two_lists(n1.next, Some(n2)) }))
            } else {
                Some(Box::new(ListNode { val: n2.val, next: merge_two_lists(Some(n1), n2.next) }))
            }
        (Some(n1), None) => Some(n1),
        (None, Some(n2)) => Some(n2),
        _ => None
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
            res.append(&mut TreeNode::preorder_traversal(node.borrow().left.clone()));
            res.append(&mut TreeNode::preorder_traversal(node.borrow().right.clone()));
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
    while i <= (n - (k-vec.len() as i32) + 1) as usize {
        vec.push(i as i32);
        backtrack_combine(vecs, vec, n, k, i + 1);
        vec.pop();

        i += 1;
    }
}

#[cfg(test)]
mod tests {
    use std::borrow::Borrow;
    use super::*;

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
                next: Some(Box::new(ListNode {
                    val: 4,
                    next: None,
                })),
            })),
        }));
        let mut l2 = Some(Box::new(ListNode {
            val: 1,
            next: Some(Box::new(ListNode {
                val: 3,
                next: Some(Box::new(ListNode {
                    val: 4,
                    next: None,
                })),
            })),
        }));
        let new_head = merge_two_lists(l1, l2);
        assert_eq!(new_head.as_ref().unwrap().val, 1);
        assert_eq!(new_head.as_ref().unwrap().next.as_ref().unwrap().val, 1);
        assert_eq!(new_head.as_ref().unwrap().next.as_ref().unwrap().next.as_ref().unwrap().val, 2);
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
                        next: Some(Box::new(ListNode {
                            val: 5,
                            next: None,
                        })),
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
        let mut root = Option::from(
            Rc::new(RefCell::new(TreeNode {
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
            })
        ));
        let result = TreeNode::preorder_traversal(root);
        assert_eq!(result, vec![1, 2, 3]);
    }

    #[test]
    fn test_level_order() {
        let mut root = Option::from(
            Rc::new(RefCell::new(TreeNode {
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
            })
        ));
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
        assert_eq!(result, vec!["((()))", "(()())", "(())()", "()(())", "()()()"]);
    }

    #[test]
    fn test_subsets() {
        let result = subsets(vec![1, 2, 3]);
        assert_eq!(result, vec![vec![], vec![1], vec![1, 2], vec![1, 2, 3], vec![1, 3], vec![2], vec![2, 3], vec![3]]);
    }

    #[test]
    fn test_combine() {
        let result = combine(4, 2);
        assert_eq!(result, vec![vec![1, 2], vec![1, 3], vec![1, 4], vec![2, 3], vec![2, 4], vec![3, 4]]);
    }
}
