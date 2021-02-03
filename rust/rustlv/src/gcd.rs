// use std::str::FromStr;
// use std::io::Write;

// fn main() {
//     let mut numbsers = Vec::new();
//
//     for arg in std::env::args().skip(1) {
//         numbsers.push(u64::from_str(&arg)
//             .expect("error parsing argument"));
//     }
//
//     if numbsers.len() == 0 {
//         writeln!(std::io::stderr(), "Usage: gcd NUMBER...").unwrap();
//         std::process::exit(1);
//     }
//
//     let mut d = numbsers[0];
//     for m in &numbsers[1..] {
//         d = gcd(d, *m);
//     }
//
//     println!("The greatest common divisor of {:?} is {}",
//              numbsers, d);
// }

fn gcd(mut n: u64, mut m: u64) -> u64 {
    assert!(n != 0 && m != 0);
    while m != 0 {
        if m < n {
            let t :u64 = m;
            m = n;
            n = t;
        }
        m = m % n;
    }
    n
}
