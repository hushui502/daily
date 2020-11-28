use std::option::Option::Some;

#[derive(PartialEq, Eq, Clone, Debug)]
pub struct ListNode {
    pub val: i32,
    pub next: Option<Box<ListNode>>
}

impl ListNode {
    fn new(val: i32) -> Self {
        ListNode {
            next: None,
            val
        }
    }
}

pub fn to_list(vec: Vec<i32>) -> Option<Box<ListNode>> {
    let mut current = None;
    for &v in vec.iter().rev() {
        let mut node = ListNode::new(v);
        node.next = current;
        current = Some(Box::new(node))
    }

    current
}

pub fn reverse_list(head: Option<Box<ListNode>>) -> Option<Box<ListNode>> {
    let mut prev = None;
    let mut curr = head;

    while let Some(mut boxed_node) = curr.take() {
        let next = boxed_node.next.take();
        boxed_node.next = prev;

        prev = Some(boxed_node);
        curr = next;
    }
    prev
}

pub fn has_cycle(head: Option<Box<ListNode>>) -> bool {
    let mut fast_p = &head;
    let mut slow_p = &head;

    while fast_p.is_some() && fast_p.as_ref().unwrap().next.is_some() {
        slow_p = &slow_p.as_ref().unwrap().next;
        fast_p= &fast_p.as_ref().unwrap().next.as_ref().unwrap().next;

        if slow_p == fast_p { return true }
    }

    false
}

pub fn merge_two_lists(l1: Option<Box<ListNode>>, l2: Option<Box<ListNode>>) -> Option<Box<ListNode>> {
    match (l1, l2) {
        (Some(node1), None) => Some(node1),
        (None, Some(node2)) => Some(node2),
        (Some(mut node1), Some(mut node2)) => {
            if node1.val < node2.val {
                let n = node1.next.take();
                node1.next = merge_two_lists(n, Some(node2));
                Some(node1)
            } else {
                let n = node2.next.take();
                node2.next = merge_two_lists(Some(node1), n);
                Some(node2)
            }
        },
        _ => None
    }
}

// pub fn middle_node(head: Option<Box<ListNode>>) -> Option<Box<ListNode>> {
//     let mut fast_p = &head;
//     let mut slow_p = &head;
//
//     while fast_p.is_some() && fast_p.as_ref().unwrap().next.is_some() {
//         slow_p = &slow_p.as_ref().unwrap().next;
//         fast_p = &fast_p.as_ref().unwrap().next.as_ref().unwrap().next;
//     }
//     slow_p.clone()
// }

pub fn remove_nth_from_end(head: Option<Box<ListNode>>, n: i32) -> Option<Box<ListNode>> {
    let mut dummy = Some(Box::new(ListNode { val: 0, next: head }));
    let mut cur = &mut dummy;
    let mut length = 0;

    while let Some(_node) = cur.as_mut() {
        cur = &mut cur.as_mut().unwrap().next;
        if let Some(_inner_node) = cur {
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

