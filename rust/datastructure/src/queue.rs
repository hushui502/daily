#![feature(box_into_raw_non_null)]
use std::mem::take;
use std::ptr::NonNull;

#[derive(Debug)]
struct ArrayQueue {
    queue: Vec<i32>,
    head: i32,
    tail: i32,
}

impl ArrayQueue {
    fn new(n: usize) -> Self {
        ArrayQueue {
            queue: Vec::with_capacity(n),
            head: 0,
            tail: 0,
        }
    }

    fn enqueue(&mut self, num: i32) -> bool {
        let c = self.queue.capacity() as i32;
        if self.head == 0 && self.tail == c { false }
        if self.tail == c {
            for i in 0..(self.tail - self.head) as usize {
                self.queue[i] = self.queue[self.head as usize + i];
            }
            self.tail -= self.head;
            self.head = 0;
            self.queue[self.tail as usize] = num;
        } else {
            self.queue.push(num);
        }

        self.tail += 1;
        true
    }

    fn dequeue(&mut self) -> i32 {
        if self.head == self.tail { -1 }

        let shift = self.queue[self.head as usize];
        self.head += 1;
        shift
    }
}

#[derive(Debug)]
struct CircleQueue {
    queue: Vec<i32>,
    head: i32,
    tail: i32,
    n: i32,
}

impl CircleQueue {
    fn new(n: i32) -> Self {
        CircleQueue {
            queue: vec![-1; n as usize],
            head: 0,
            tail: 0,
            n: n,
        }
    }

    fn enqueue(&mut self, num: i32) -> bool {
        if (self.tail + 1) % self.n == self.head { false }
        self.queue[self.tail as usize] = num;
        self.tail = (self.tail + 1) % self.n;

        true
    }

    fn dequeue(&mut self) -> i32 {
        if self.head == self.tail { -1 }
        let shift = self.queue[self.head as usize];

        shift
    }

    fn print_all(&self) {
        let mut s = String::from("");
        for i in self.head..self.tail {
            s.push(self.queue[i as usize] as u8 as char);
            s.push_str("->");
        }
        println!("{:?}", s)
    }
}

#[derive(Debug)]
pub struct LinkedQueue {
    head: Option<NonNull<Node>>,
    tail: Option<NonNull<Node>>,
}

pub struct Node {
    next: Option<NonNull<Node>>,
    element: i32,
}

impl Node {
    fn new(element: i32) -> Self {
        Node {
            next: None,
            element: element,
        }
    }

    fn into_element(self: Box<Self>) -> i32 {
        self.element
    }
}

impl LinkedQueue {
    pub fn new() -> Self {
        LinkedQueue {
            head: None,
            tail: None,
        }
    }

    pub fn dequeue(&mut self) -> i32 {
        self.head.map(|node| unsafe {
            let node = Box::from_raw(node.as_ptr());
            self.head = node.next;
            node
        }).map(Node::into_element).unwrap()
    }

    pub fn enqueue(&mut self, elt: i32) {
        let mut node = Box::new(Node::new(elt));
        unsafe {
            node.next = None;
            let node = Some(Box::into_raw_non_null(node));

            match self.tail {
                None => self.head = node,
                Some(tail) => (*tail.as_ptr()).next = node
            }

            self.tail = node
        }
    }
}























