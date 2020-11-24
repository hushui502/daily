use std::collections::HashMap;

fn main() {
    let mut v = vec![1, 2, 3];
    let third: &i32 = &v[2];
    println!("The third element is {}", third);

    match v.get(0) {
        Some(x) => println!("The third element is {}", x),
        None => println!("There is no third element.")
    }

    for i in &mut v {
        *i += 10;
        println!("{}", i)
    }


    enum SpreadsheetCell {
        Int(i32),
        Float(f64),
        Text(String),
    }

    let row = vec![
        SpreadsheetCell::Int(3),
        SpreadsheetCell::Text(String::from("blue")),
        SpreadsheetCell::Float(10.22),
    ];


    let text = "hello world   beijing hufan ";
    let mut map = HashMap::new();

    for word in text.split_whitespace() {
        let count = map.entry(word).or_insert(0);
        *count += 1;
    }

    println!("{:?}", map);

    panic!("cash down")

}


