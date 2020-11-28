
fn main() {

    let v = [10, 40, 30];
    assert_eq!(Some(&30), v.last());

    let w: &[i32] = &[];
    assert_eq!(None, w.last());

    println!("{:?}", v)
}

#[test]
fn test_last() {
    let v = [10, 40, 30];
    assert_eq!(Some(&30), v.last());

    let w: &[i32] = &[];
    assert_eq!(None, w.last());

    println!("{:?}", v)
}
