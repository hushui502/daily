use std::i32::MAX;
use std::mem::swap;
use std::rc::Rc;
use std::collections::HashMap;

fn main() {
    // let big_val = MAX;
    // let d = big_val + 1;

    // 反转wrap为负值
    // let x = big_val.wrapping_add(1);
    // println!("{}", x);

    // assert_eq!(2u16.pow(4), 16);
    // assert_eq!((-4i32).abs(), 4)
    // assert_eq!(0b10101u8.count_ones(), 3)

    // assert_eq!(5f32.sqrt() * 5f32.sqrt(), 5.);
    // assert_eq!(-1.22f64.floor(), -1.0);
    // assert!((-1. / std::f32::INFINITY).is_sign_negative());
    //
    // println!("{}", (2.2f64).sqrt());
    // println!("{}", f64::sqrt(2.2));
    //
    // assert_eq!(false as i32, 0);
    // assert_eq!(true as i32, 1);

    // println!("\u{CA0}");
    // assert_eq!(std::char::from_digit(2, 10), Some('2'));
    //
    // let text = "hello world";
    // let (head, tail) = text.split_at(text.len()/2);
    // let temp = text.split_at(3);
    // println!("{} {}", head, tail);
    // println!("{} {}", temp.0, temp.1);
    //
    // let mut a = 1;
    // let mut b = 2;
    // swap(&mut a, &mut b);
    // println!("{}", a);
    //
    // let t = (12, "aaa");
    // let b = Box::new(t);
    //
    // let mut sieve = [true, 1000];
    // for i in 1..100 {
    //     if sieve[i] {
    //         let mut j = i * i;
    //         while j < 1000 {
    //             sieve[j] = false;
    //             j += i;
    //         }
    //     }
    // }

    // let mut chaos = [1, 3, 2, 5, 4];
    // chaos.sort();
    // assert_eq!(chaos, [1, 2, 3, 4, 5]);
    //
    // let mut v = vec![1, 2, 3];
    // assert_eq!(v.iter().fold(1, |a, b| a * b), 6);

    // let mut vv = Vec::new();
    // vv.push(11);
    // vv.push(22);
    // assert_eq!(vv, vec![11, 22]);

    // let v: Vec<i32> = (0..5).collect();
    // assert_eq!(v, [0, 1, 2, 3, 4]);

    // let mut v = vec!["a", "b", "c"];
    // v.reverse();
    // assert_eq!(v, vec!["c", "b", "a"]);

    // let mut v = Vec::with_capacity(2);
    // assert_eq!(v.len(), 0);
    // assert_eq!(v.capacity(), 2);
    // v.push(1);
    // v.push(2);
    // v.push(2);
    // // 4
    // assert_eq!(v.capacity(), 2);
    // v.insert(3, 222);
    // v.remove(4);
    // assert_eq!(v.pop(), Some("aa"));
    //
    // let languages: Vec<String> = std::env::args().skip(1).collect();
    // for l in languages {
    //     println!("{}: {}", l,
    //              if l.len() % 2 == 0 {
    //                  "functional"
    //              } else {
    //                  "imperative"
    //              });
    // }

    // let v: Vec<f64> = vec![1.1, 2.1, 3.1];
    // let a: [f64; 3] = [1.1, 2.2, 3.3];
    //
    // let sv = &v;
    // let sa = &a;
    //
    // print(&v[0..1]);
    // print(&a[1..2]);

    // let bits = vec!["aa", "bb", "cc"];
    // assert_eq!(bits.concat(), "aabbcc");
    // assert_eq!(bits.join(", "));

    // assert_eq!("Hell".to_lowercase(), "hell");
    // assert!("helL".to_lowercase().contains("ll"));
    // assert_eq!(" clean  up".trim(), "clean  up");
    //
    // for word in "aa, ab, ac".split(", ") {
    //     assert!(word.starts_with("v"))
    // }

    // {
    //     let point = Box::new((1, 2, 3));
    //     let label = format!("{:?}", point);
    //     println!("{:?}", point);
    //     println!("{:?}", label);
    // }

    // struct Person { name:String, birth: i32 }
    //
    // let mut composers = Vec::new();
    // composers.push(Person { name: "hufan1".to_string(), birth: 123 });
    // composers.push(Person { name: "hufan2".to_string(), birth: 223 });
    // composers.push(Person { name: "hufan3".to_string(), birth: 333 });
    //
    // for composer in &composers {
    //     println!("{}, born {}", composer.name, composer.birth);
    // }
    // for composer in &composers {
    //     println!("{}, born {}", composer.name, composer.birth);
    // }

    // let s = vec!["hello".to_string(), "world".to_string()];
    // let t = s.clone();
    // let u = s.clone();

    // let mut s = "hello".to_string();
    // let u = s;
    // s = "hello".to_string();

    // let mut v = Vec::new();
    // for i in 101 .. 106 {
    //     v.push(i.to_string());
    // }
    //
    // let fifith = v.pop().unwrap();
    // assert_eq!(fifith, "105");
    // assert_eq!(v.len(), 4);
    //
    // let second = v.swap_remove(1);
    // assert_eq!(second, "102");
    // assert_eq!(v, ["101", "104", "103"]);
    //
    // let third = std::mem::replace(&mut v[2], "substitute".to_string());
    // assert_eq!(third, "103");
    // assert_eq!(v, vec!["101", "104", "substitute"]);

    // let v = vec!["aa".to_string(), "bb".to_string(), "cc".to_string()];
    // for mut s in v {
    //     s.push('!');
    //     println!("{}", s);
    // }
    // println!("{:?}", v)

    // struct Person { name: Option<String>, birth: i32 }
    // let mut composers = Vec::new();
    // composers.push(Person { name: Some("hufan".to_string()), birth: 1234 });
    //
    // // let first_name = composers[0].name;      // con not move, because name maybe None
    // let first_name = composers[0].name.take();
    // // let first_name = std::mem::replace(&mut composers[0].name, None);
    // assert_eq!(first_name, Some("hufan".to_string()));
    // assert_eq!(composers[0].name, None);


    // let num1: i32 = 12;
    // let num2 = num1;
    // println!("{}", num1);

    /// 只要copy clone过的所有字段都是Copy类型的对象都可以实现c++种的拷贝复制了，不再像Rust的转移了，sad
    // #[derive(Copy, Clone)]
    // struct Label { number: u32 }
    // fn print(l: Label) { println!("{:?}", l.number)}
    // let l = Label { number: 3 };
    // print(l);
    // println!("{}", l.number);
    //
    // #[derive(Copy, Clone)]
    // struct StringLabel { name: String }
    // fn printStringLabel(l: StringLabel) { println!("{}", l.name)}
    // let strL = StringLabel { name: "hello".to_string() };
    // printStringLabel(strL);
    // println!("{}", strL.name);

    // let s: Rc<String> = Rc::new("shirnk".to_string());
    // let t = s.clone();
    // let u = s.clone();
    // assert!(s.contains("shirnk"));
    // assert_eq!(t.find("irnk"), Some(2));
    // println!("{}", u);
    // println!("{}", s);
    // println!("{}", s);
    // println!("{}", u);
    // /// 不允许y引用指针修改原引用，因为是会导致并发下的数据不一致
    // s.push_str("node");

    type Table = HashMap<String, Vec<String>>;

    let mut table = Table::new();
    table.insert("hufan".to_string(),
                 vec!["a2".to_string(), "a1".to_string()]);
    table.insert("libai".to_string(),
                 vec!["l1".to_string(), "l2".to_string()]);

    sort_works(&mut table);
    show(&table);

    // show(& mut table);
    fn show(table: &Table) {
        for (artist, works) in table {
            println!("works by {}", artist);
            for work in works {
                println!("{}", work);
            }
        }
    }

    fn sort_works(table: &mut Table) {
        for (_artist, works) in table {
            works.sort()
        }
    }
}



fn print(n: &[f64]) {
    for ele in n {
        println!("{}", ele);
    }
}