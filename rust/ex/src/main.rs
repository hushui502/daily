use std::mem;
use std::convert::{TryFrom, TryInto};

static LANGUAGE: &'static str = "Rust";

const THRESHOLD: i32 = 10;

fn is_big(n: i32) -> bool {
    n > THRESHOLD
}

type U8 = u8;


#[derive(Debug)]
struct Number {
    value: i32,
}

impl From<i32> for Number {
    fn from(item: i32) -> Self {
        Number { value: item }
    }
}

#[derive(Debug, PartialEq)]
struct EvenNumber(i32);

impl TryFrom<i32> for EvenNumber {
    type Error = ();

    fn try_from(value: i32) -> Result<Self, Self::Error> {
        if value % 2 == 0 {
            Ok(EvenNumber(value))
        } else {
            Err(())
        }
    }
}

struct Circle {
    radius: i32
}

impl ToString for Circle {
    fn to_string(&self) -> String {
        format!("Circle of radius is {:?}", self.radius)
    }
}


fn main() {
    // assert_eq!(EvenNumber::try_from(9), Ok(EvenNumber(9)));
    assert_eq!(EvenNumber::try_from(8), Ok(EvenNumber(8)));


    let result: Result<EvenNumber, ()> = 8i32.try_into();
    assert_eq!(result, Ok(EvenNumber(8)));

    let result: Result<EvenNumber, ()> = 5i32.try_into();
    assert_eq!(result, Err(()));

    let circle = Circle { radius: 4 };
    println!("{}", circle.to_string());


    let parsed: i32 = "4".parse().unwrap();
    let turbo_parsed = "10".parse::<i32>().unwrap();

    let sum = parsed + turbo_parsed;
    println!("{:?}", sum);

    let mut count = 1;
    loop {
        count += 1;
        println!("=== {}", count);
        if count == 10 {
            break
        }
    }

    'outer: loop {
        println!("Entered the outer loop");
        'linner: loop {
            println!("Entered the inner loop");
            break 'outer
        }
    }

    let res = loop {
        count += 1;
        if count == 1000 {
            break count * 2
        }
    };
    assert_eq!(res, 2000);

    for i in 1..=10 {
        println!("{}", i)
    }

    let mut names = vec!["hufan", "libai", "zhangfei"];
    for name in names.iter_mut() {
        *name = match name {
            &mut "hufan" => "hufahufahufauhu",
            _ => "hello"
        }
    }
    println!("{:?}", names);


    let number = 3;
    match number {
        1 => println!("one"),
        2 | 3 | 4 => println!("two --- four"),
        8..=10 => println!("eight --- ten"),
        _ => println!("> 10")
    }


}