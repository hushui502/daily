use std::collections::{HashMap, VecDeque};

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
    val: i32,
    next: Option<Box<ListNode>>,
}

pub fn reverse_list(head: Option<Box<ListNode>>) -> Option<Box<ListNode>> {
    if head.is_none() {
        return None;
    }
    let mut prev: Option<Box<ListNode>> = None;
    let mut curr = head;
    while let Some(mut node) = curr {
        let next = node.next.take();
        node.next = prev.take();
        prev = Some(node);
        curr = next;
    }

    prev
}

#[cfg(test)]
mod tests {
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

    #[test]
    fn test_reverse_list() {
        let mut head = Some(Box::new(ListNode {
            val: 1,
            next: Some(Box::new(ListNode {
                val: 2,
                next: Some(Box::new(ListNode {
                    val: 3,
                    next: Some(Box::new(ListNode { val: 4, next: None })),
                })),
            })),
        }));
        let mut new_head = reverse_list(head);
        assert_eq!(new_head.as_ref().unwrap().val, 4);
        assert_eq!(new_head.unwrap().next.unwrap().val, 3);
    }
}
