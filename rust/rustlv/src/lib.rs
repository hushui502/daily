#![allow(dead_code)]
//                                  my_ascii
// mod my_ascii {
//     #[derive(Debug, Eq, PartialEq)]
//     pub struct Ascii (
//         Vec<u8>
//     );
//
//     impl Ascii {
//         pub fn from_bytes(bytes: Vec<u8>) -> Result<Ascii, NotAsciiError> {
//             if bytes.iter().any(|&byte| !byte.is_ascii()) {
//                 return Err(NotAsciiError(bytes))
//             }
//             Ok(Ascii(bytes))
//         }
//     }
//
//     #[derive(Debug, Eq, PartialEq)]
//     pub struct NotAsciiError(pub Vec<u8>);
//
//     impl From<Ascii> for String {
//         fn from(ascii: Ascii) -> String {
//             unsafe {
//                 String::from_utf8_unchecked(ascii.0)
//             }
//         }
//     }
//
//     impl Ascii {
//         pub unsafe fn from_bytes_unchecked(bytes: Vec<u8>) -> Ascii {
//             Ascii(bytes)
//         }
//     }
// }
//
// #[test]
// fn good_ascii() {
//     use my_ascii::Ascii;
//
//     let bytes: Vec<u8> = b"ASCII and ye shall receive".to_vec();
//
//     let ascii: Ascii = Ascii::from_bytes(bytes)
//         .unwrap();
//
//     let string = String::from(ascii);
//
//     assert_eq!(string, "ASCII and ye shall receive");
// }
//
// fn bad_ascii() {
//     use my_ascii::Ascii;
//
//     let bytes = vec![0xf7, 0xbf, 0xbf, 0xbf];
//
//     let ascii = unsafe {
//         Ascii::from_bytes_unchecked(bytes)
//     };
//
//     let bogus: String = ascii.into();
//
//     assert_eq!(bogus.chars().next().unwrap() as u32, 0x1fffff);
// }


use std::io::empty;

//                              binary-tree
// An ordered collection of 'T's
enum BinaryTree<T> {
    Empty,
    NonEmpty(Box<TreeNode<T>>),
}

// A part of a BinaryTree.
struct TreeNode<T> {
    element: T,
    left: BinaryTree<T>,
    right: BinaryTree<T>,
}

impl <T: Clone> BinaryTree<T> {
    fn walk(&self) -> Vec<T> {
        match *self {
            BinaryTree::Empty => vec![],
            BinaryTree::NonEmpty(ref boxed) => {
                let mut result = boxed.left.walk();
                result.push(boxed.element.clone());
                result.extend(boxed.right.walk());
                result
            }
        }
    }
}

impl<T: Ord> BinaryTree<T> {
    fn add(&mut self, value: T) {
        match *self {
            BinaryTree::Empty => {
                *self = BinaryTree::NonEmpty(Box::new(TreeNode {
                    element: value,
                    left: BinaryTree::Empty,
                    right: BinaryTree::Empty,
                }))
            }
            BinaryTree::NonEmpty(ref mut node) => {
                if value <= node.element {
                    node.left.add(value);
                } else {
                    node.right.add(value);
                }
            }
        }
    }
}

#[test]
fn binary_tree_size() {
    use std::mem::size_of;

    let word = size_of::<usize>();
    assert_eq!(size_of::<BinaryTree<String>>(), word);

    type Triple = (&'static str, BinaryTree<&'static str>, BinaryTree<&'static str>);
    assert_eq!(size_of::<Triple>(), 4 * word);
}

#[test]
fn build_binary_tree() {
    use self::BinaryTree::*;

    let jupiter_tree = NonEmpty(Box::new(TreeNode {
        element: "Jupiter",
        left: Empty,
        right: Empty,
    }));

    let mercury_tree = NonEmpty(Box::new(TreeNode {
        element: "Mercury",
        left: Empty,
        right: Empty,
    }));

    let mars_tree = NonEmpty(Box::new(TreeNode {
        element: "Mars",
        left: jupiter_tree,
        right: mercury_tree,
    }));

    let venus_tree = NonEmpty(Box::new(TreeNode {
        element: "Venus",
        left: Empty,
        right: Empty,
    }));

    let uranus_tree = NonEmpty(Box::new(TreeNode {
        element: "Uranus",
        left: Empty,
        right: venus_tree,
    }));

    let tree = NonEmpty(Box::new(TreeNode {
        element: "Saturn",
        left: mars_tree,
        right: uranus_tree,
    }));

    assert_eq!(tree.walk(),
               vec!["Jupiter", "Mars", "Mercury", "Saturn", "Uranus", "Venus"]);
}

#[test]
fn test_add_method_1() {
    let planets = vec!["Mercury", "Venus", "Mars", "Jupiter", "Saturn", "Uranus"];
    let mut tree = BinaryTree::Empty;
    for planet in planets {
        tree.add(planet);
    }

    assert_eq!(tree.walk(),
               vec!["Jupiter", "Mars", "Mercury", "Saturn", "Uranus", "Venus"]);
}

#[test]
fn test_add_method_2() {
    let mut tree = BinaryTree::Empty;
    tree.add("Mercury");
    tree.add("Venus");
    for planet in vec!["Mars", "Jupiter", "Saturn", "Uranus"] {
        tree.add(planet);
    }

    assert_eq!(
        tree.walk(),
        vec!["Jupiter", "Mars", "Mercury", "Saturn", "Uranus", "Venus"]
    );
}


// Iterators
use self::BinaryTree::*;

struct TreeIter<'a, T> {
    unvisited: Vec<&'a TreeNode<T>>
}

impl<'a, T: 'a> TreeIter<'a, T> {
    fn push_left_edge(&mut self, mut tree: &'a BinaryTree<T>) {
        while let NonEmpty(ref node) = *tree {
            self.unvisited.push(node);
            tree = &node.left;
        }
    }
}

impl<T> BinaryTree<T> {
    fn iter(&self) -> TreeIter<T> {
        let mut iter = TreeIter { unvisited: Vec::new() };
        iter.push_left_edge(self);
        iter
    }
}

impl<'a, T: 'a> IntoIterator for &'a BinaryTree<T> {
    type Item = &'a T;
    type IntoIter = TreeIter<'a, T>;
    fn into_iter(self) -> Self::IntoIter {
        self.iter()
    }
}

impl<'a, T> Iterator for TreeIter<'a, T> {
    type Item = &'a T;
    fn next(&mut self) -> Option<&'a T> {
        let node = self.unvisited.pop()?;
        self.push_left_edge(&node.right);

        Some(&node.element)
    }
}

#[test]
fn external_iterator() {
    fn make_node<T>(left: BinaryTree<T>, element: T, right: BinaryTree<T>) -> BinaryTree<T> {
        NonEmpty(Box::new(TreeNode { left, element, right}))
    }

    let mut tree = BinaryTree::Empty;
    tree.add("jaeger");
    tree.add("robot");
    tree.add("droid");
    tree.add("mecha");

    let mut v = Vec::new();
    for kind in &tree {
        v.push(*kind);
    }
    assert_eq!(v, ["droid", "jaeger", "mecha", "robot"]);

    assert_eq!(tree.iter()
        .map(|name| format!("mega-{}", name))
        .collect::<Vec<_>>(),
               vec!["mega-droid", "mega-jaeger", "mega-mecha", "mega-robot"]);

    let left_subtree = make_node(Empty, "mecha", Empty);
    let right_subtree = make_node(make_node(Empty, "droid", Empty),
    "robot", Empty);
    let tree = make_node(left_subtree, "Jaeger", right_subtree);

    let mut v = Vec::new();
    let mut iter = TreeIter { unvisited: vec![] };
    iter.push_left_edge(&tree);
    for kind in iter {
        v.push(*kind);
    }
    assert_eq!(v, ["mecha", "Jaeger", "droid", "robot"]);

    let mut v = Vec::new();
    for kind in &tree {
        v.push(*kind);
    }
    assert_eq!(v, ["mecha", "Jaeger", "droid", "robot"]);
}

#[test]
fn other_cloned() {
    use std::collections::BTreeSet;

    let mut set = BTreeSet::new();
    set.insert("mecha");
    set.insert("Jaeger");
    set.insert("droid");
    set.insert("robot");
    assert_eq!(set.iter().cloned().collect::<Vec<_>>(),
               ["Jaeger", "droid", "mecha", "robot"]);
}

#[test]
fn fuzz() {
    fn make_random_tree(p: f32) -> BinaryTree<i32> {
        use rand::prelude::*;
        use rand::thread_rng;
        use rand::rngs::ThreadRng;

        fn make(p: f32, next: &mut i32, rng: &mut ThreadRng) -> BinaryTree<i32> {
            if rng.gen_range(0.0 .. 1.0) > p {
                Empty
            } else {
                let left = make(p * p, next, rng);
                let element = *next;
                *next += 1;
                let right = make(p * p, next, rng);
                NonEmpty(Box::new(TreeNode { left, element, right }))
            }
        }

        make(p, &mut 0, &mut thread_rng())
    }

    for _ in 0..100 {
        let tree = make_random_tree(0.9999);
        assert!(tree.into_iter().fold(Some(0), |s, &i| {
            s.and_then(|expected| if i == expected { Some(expected+1) } else { None })
        }).is_some());
    }
}

//                              Complex number examples.

macro_rules! define_complex {
    () => {
        #[derive(Clone, Copy, Debug)]
        struct Complex<T> {
            /// Real portion of the complex number
            re: T,

            /// Imaginary portion of the complex number
            im: T,
        }
    };
}

mod first_cut {
    #[derive(Clone, Copy, Debug)]
    struct Complex<T> {
        /// Real portion of the complex number
        re: T,

        /// Imaginary portion of the complex number
        im: T,
    }

    use std::ops::Add;

    impl<T> Add for Complex<T>
        where
            T: Add<Output = T>,
    {
        type Output = Self;
        fn add(self, rhs: Self) -> Self {
            Complex {
                re: self.re + rhs.re,
                im: self.im + rhs.im,
            }
        }
    }

    use std::ops::Sub;

    impl<T> Sub for Complex<T>
        where
            T: Sub<Output = T>,
    {
        type Output = Self;
        fn sub(self, rhs: Self) -> Self {
            Complex {
                re: self.re - rhs.re,
                im: self.im - rhs.im,
            }
        }
    }

    use std::ops::Mul;

    impl<T> Mul for Complex<T>
        where
            T: Clone + Add<Output = T> + Sub<Output = T> + Mul<Output = T>,
    {
        type Output = Self;
        fn mul(self, rhs: Self) -> Self {
            Complex {
                re: self.re.clone() * rhs.re.clone()
                    - (self.im.clone() * rhs.im.clone()),
                im: self.im * rhs.re + self.re * rhs.im,
            }
        }
    }

    #[test]
    fn try_it_out() {
        let mut z = Complex { re: 1, im: 2 };
        let c = Complex { re: 3, im: 4 };

        z = z * z + c;

        std::mem::forget(z);
    }

    impl<T: PartialEq> PartialEq for Complex<T> {
        fn eq(&self, other: &Complex<T>) -> bool {
            self.re == other.re && self.im == other.im
        }
    }

    #[test]
    fn test_complex_eq() {
        let x = Complex { re: 5, im: 2 };
        let y = Complex { re: 2, im: 5 };
        assert_eq!(x * y, Complex { re: 0, im: 29 });
    }

    // impl<T: Eq> Eq for Complex<T> {}
}

mod non_generic_add {
    define_complex!();

    use std::ops::Add;

    impl Add for Complex<i32> {
        type Output = Complex<i32>;

        fn add(self, rhs: Self) -> Self {
            Complex {
                re: self.re + rhs.re,
                im: self.im + rhs.im,
            }
        }
    }
}

mod somewhat_generic {
    define_complex!();

    use std::ops::Add;

    impl<T> Add for Complex<T>
        where
            T: Add<Output = T>,
    {
        type Output = Self;
        fn add(self, rhs: Self) -> Self {
            Complex {
                re: self.re + rhs.re,
                im: self.im + rhs.im,
            }
        }
    }

    use std::ops::Neg;

    impl<T> Neg for Complex<T>
    where
        T: Neg<Output = T>,
    {
        type Output = Complex<T>;

        fn neg(self) -> Complex<T> {
            Complex {
                re: -self.re,
                im: -self.im,
            }
        }
    }
}

mod very_generic {
    define_complex!();

    use std::ops::Add;
    impl<L, R> Add<Complex<R>> for Complex<L>
    where
        L: Add<R>,
    {
        type Output = Complex<L::Output>;

        fn add(self, rhs: Complex<R>) -> Self::Output {
            Complex {
                re: self.re + rhs.re,
                im: self.im + rhs.im,
            }
        }
    }
}

mod impl_compound {
    define_complex!();

    use std::ops::AddAssign;

    impl<T> AddAssign for Complex<T>
        where
            T: AddAssign<T>,
    {
        fn add_assign(&mut self, rhs: Complex<T>) {
            self.re += rhs.re;
            self.im += rhs.im;
        }
    }
}

mod derive_partialeq {
    #[derive(Clone, Copy, Debug, PartialEq)]
    struct Complex<T> {
        re: T,
        im: T,
    }
}

mod derive_everything {
    #[derive(Clone, Copy, Debug, Eq, PartialEq)]
    struct Complex<T> {
        /// Real portion of the complex number
        re: T,

        /// Imaginary portion of the complex number
        im: T,
    }
}

mod formatting {
    use std::fmt::Formatter;

    #[test]
    fn complex() {
        #[derive(Copy, Clone, Debug)]
        struct Complex { re: f64, im: f64 }

        let third = Complex { re: -0.5, im: f64::sqrt(0.75) };
        println!("{:?}", third);

        use std::fmt;

        impl fmt::Display for Complex {
            fn fmt(&self, dest: &mut fmt::Formatter) -> fmt::Result {
                let im_sign = if self.im < 0.0 { '-' } else { '+' };
                write!(dest, "{} {} {}i", self.re, im_sign, f64::abs(self.im))
            }
        }

        let one_twenty = Complex { re: -0.5, im: 0.866 };
        assert_eq!(format!("{}", one_twenty),
                   "-0.5 + 0.866i");

        let two_forty = Complex { re: -0.5, im: -0.866 };
        assert_eq!(format!("{}", two_forty),
                   "-0.5 - 0.866i");
    }

    #[test]
    fn complex_fancy() {
        #[derive(Copy, Clone, Debug)]
        struct Complex { re: f64, im: f64 }

        use std::fmt;

        impl fmt::Display for Complex {
            fn fmt(&self, dest: &mut fmt::Formatter) -> fmt::Result {
                let (re, im) = (self.re, self.im);
                if dest.alternate() {
                    let abs = f64::sqrt(re * re + im * im);
                    let angle = f64::atan2(im, re) / std::f64::consts::PI * 180.0;
                    write!(dest, "{} ∠ {}°", abs, angle)
                } else {
                    let im_sign = if im < 0.0 { '-' } else { '+' };
                    write!(dest, "{} {} {}i", re, im_sign, f64::abs(im))
                }
            }
        }

        let ninety = Complex { re: 0.0, im: 2.0 };
        assert_eq!(format!("{}", ninety),
                   "0 + 2i");
        assert_eq!(format!("{:#}", ninety),
                   "2 ∠ 90°");
    }
}

