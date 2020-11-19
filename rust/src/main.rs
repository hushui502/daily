// use std::io;
// use std::cmp::Ordering;
// use rand::Rng;
use std::net::{IpAddr, Ipv4Addr, Ipv6Addr};

use std::mem::take;

fn main() {
    let localhost_v4 = IpAddr::V4(Ipv4Addr::new(127, 0, 0, 1));
    let localhost_v6 = IpAddr::V6(Ipv6Addr::new(0, 0, 0, 0, 0, 0, 0, 1));

    assert_eq!("127.0.0.1".parse(), Ok(localhost_v4));
    assert_eq!("::1".parse(), Ok(localhost_v6));

    assert_eq!(localhost_v4.is_ipv6(), false);
    assert_eq!(localhost_v4.is_ipv4(), true);

    let m = Message::Write(String::from("hello"));
    m.call();

    let mut x = Some(2);
    let mut y = Some(3);
    match x.as_mut() {
        Some(v) => *v = 42,
        None => {}
    }
    assert_eq!(x, Some(42));
    let sum = x.and(y);
    // println!("{}", sum.)

    let penny = Coin::Penny;
    let n = value_in_cents(penny);


    // vector
    let v = vec![1, 2, 3];
    let third: &i32 = &v[2];
    println!("third element is {}", &third)


}


enum Coin {
    Penny,
    Nickel,
    Dime,
    Quarter,
}

fn value_in_cents(coin: Coin) -> u8 {
    match coin {
        Coin::Penny => {
            println!("{}", 1);
            1
        },
        Coin::Quarter => 4,
        _ => -1
    }

    let mut count = 0;
    if let Coin::Quarter = coin {
        4
    } else {
        count += 1;
    }
}

enum Message {
    Quit,
    Move{x: i32, y: i32},
    Write(String),
    ChangeColor(i32, i32, i32)
}

impl Message {
    fn call(&self) {

    }
}

struct QuitMessage;
struct MoveMessage {
    x: i32,
    y: i32,
}
struct WriteMessage(String);
struct ChangeColorMessage(i32, i32, i32);

enum IpAddrKind {
    V4(u8, u8, u8, u8),
    V6,
}

// struct IpAddr {
//     kind: IpAddrKind,
//     address: String,
// }


#[derive(Debug)]
struct Rectangle {
    width: u32,
    height: u32,
}

impl Rectangle {
    fn area(&self) -> u32 {
        self.width * self.height
    }
    fn can_hold(&self, other: &Rectangle) -> bool {
        self.width > other.width && self.height > other.height
    }
    fn square(size: u32) -> Rectangle {
        Rectangle {width: size, height: size}
    }
}


struct Color(i32, i32, i32);

struct User {
    username: String,
    email: String,
    sign_in_count: u64,
    active: bool,
}

fn build_user(email: String, username: String) -> User {
    User {
        email,
        username,
        active: true,
        sign_in_count: 4
    }
}

fn first_word(s: &String) -> &str {
    let bytes = s.as_bytes();

    for (i, &item) in bytes.iter().enumerate() {
        if item == b' ' {
            return &s[0..i];
        }
    }

    &s[..]
}

// fn dangle() ->&String {
//     let s = String::from("hello");
//
//     &s
// }

fn change(some_string: &mut String) {
    some_string.push_str(" world");
}

fn calculate_len(s: &String) -> (&String, usize) {
    let len = s.len();
    (s, len)
}

fn take_str(x: String) {
    println!("{}", x)
}

fn another_method(x: i32) {
    println!("The value is {}", x);
}

fn five() -> i32 {
    5
}

fn gives_str() -> String {
    let some_str = String::from("hello");
    some_str
}


