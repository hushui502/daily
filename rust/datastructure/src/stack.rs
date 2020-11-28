use std::thread::sleep;

#[derive(Hash, Eq, PartialEq, Debug, Default, Clone)]
pub struct ListNode {
    val: String,
    next: Option<Box<ListNode>>,
}

#[derive(Hash, Eq, PartialEq, Debug, Default, Clone)]
pub struct LinkedListStack {
    node: Option<Box<ListNode>>,
}

impl ListNode {
    fn new(val: String) -> Self {
        ListNode { val: val, next: None }
    }
}

impl LinkedListStack {
    pub fn new() -> Self {
        Default::default()
    }

    pub fn push(&mut self, x: String) {
        let mut n = ListNode::new(x);
        n.next = self.node.clone();
        self.node = Some(Box::new(n));
    }

    pub fn pop(&mut self) -> String {
        if self.is_empty() { return "-1".to_string(); }

        let val = self.node.as_ref().unwrap().val.clone();
        self.node = self.node.as_mut().unwrap().next.take();
        val.to_string()
    }

    pub fn print_all(&mut self) {
        let mut list = String::from("");

        while let Some(n) = self.node.as_mut() {
            list.push_str(&(n.val).to_string());
            list.push_str("-->");
            self.node = n.next.take();
        }
        println!("{:?}", list);
    }

    pub fn clear(&mut self) {
        self.node = None;
    }

    pub fn is_empty(&self) -> bool {
        if self.node.is_none() { true } else { false }
    }
}
#[derive(Debug)]
struct ArrayStack {
    data: Vec<i32>,
    top: i32,
}

impl ArrayStack {
    fn new() -> Self {
        ArrayStack { data: Vec::with_capacity(32), top: -1 }
    }

    fn push(&mut self, x: i32) {
        self.top += 1;
        if self.top == self.data.capacity() as i32 {
            let tmp_arr = self.data.clone();
            self.data = Vec::with_capacity(self.data.capacity() * 2);
            for d in tmp_arr.into_iter() {
                self.data.push(d);
            }
        }
        self.data.push(x);
    }

    fn pop(&mut self) {
        if self.is_empty() { return; }
        self.top -= 1;
        self.data.pop();
    }

    fn top(&self) -> i32 {
        if self.is_empty() { return -1; }
        *self.data.last().unwrap()
    }

    fn is_empty(&self) -> bool {
        if self.top < 0 { true } else { false }
    }
}

#[derive(Hash, Eq, PartialEq, Debug, Default, Clone)]
struct Browser {
    forward_stack: LinkedListStack,
    back_stack: LinkedListStack,
    current_page: String,
}

impl Browser {
    fn new() -> Self {
        Default::default()
    }

    fn open(&mut self, url: String) {
        if !self.current_page.is_empty() {
            self.back_stack.push(self.current_page.clone());
            self.forward_stack.clear();
        }
        self.show_url(&url, "Open".to_string())
    }

    fn go_back(&mut self) -> String {
        if self.can_go_back() {
            self.forward_stack.push(self.current_page.clone());
            let back_url = self.back_stack.pop();
            self.show_url(&back_url, "Back".to_string());

            back_url;
        }

        println!("Can not go back!");
        "-1".to_string()
    }

    fn go_forward(&mut self) -> String {
        if self.can_go_forword() {
            self.back_stack.push(self.current_page.clone());
            let forward_url = self.forward_stack.pop();
            self.show_url(&forward_url, "Forward".to_string())

            forward_url
        }

        println!("Can not go forward");
        "-1".to_string()
    }

    fn can_go_forword(&self) -> bool {
        !self.forward_stack.is_empty()
    }

    fn can_go_back(&self) -> bool {
        !self.back_stack.is_empty()
    }

    fn show_url(&mut self, url: &String, prefix: String) {
        self.current_page = url.to_string();
        println!("{:?} page == {:?}", prefix, url);
    }

    fn check_current_page(&self) {
        println!("current page == {:?}", self.current_page)
    }
}