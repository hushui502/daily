#![allow(dead_code)]

//                                  router
// use std::collections::HashMap;
//
// struct Request {
//     method: String,
//     url: String,
//     headers: HashMap<String, String>,
//     body: Vec<u8>
// }
//
// struct Response {
//     code: u32,
//     headers: HashMap<String, String>,
//     body: Vec<u8>
// }
//
// type BoxedCallback = Box<dyn Fn(&Request) -> Response>;
//
// struct BasicRouter {
//     routes: HashMap<String, BoxedCallback>
// }
//
// impl BasicRouter {
//     fn new() -> BasicRouter {
//         BasicRouter { routes: HashMap::new() }
//     }
//
//     fn add_route<C>(&mut self, url: &str, callback: C)
//         where C: Fn(&Request) -> Response + 'static
//     {
//         self.routes.insert(url.to_string(), Box::new(callback));
//     }
// }
//
// impl BasicRouter {
//     fn handle_request(&self, request: &Request) -> Response {
//         match self.routes.get(&request.url) {
//             Some(callback) => callback(request),
//             None => not_found_response()
//         }
//     }
// }
//
// fn not_found_response() -> Response {
//     Response {
//         code: 404,
//         headers: HashMap::new(),
//         body: b"<h1>Page not found</h1>".to_vec()
//     }
// }
//
// fn get_form_response() -> Response {
//     Response {
//         code: 200,
//         headers: HashMap::new(),
//         body: b"<form>".to_vec()
//     }
// }
//
// fn get_gcd_response(_req: &Request) -> Response {
//     Response {
//         code: 500,
//         headers: HashMap::new(),
//         body: b"<h1>Internal server error</h1>".to_vec()
//     }
// }
//
// fn req(url: &str) -> Request {
//     Request {
//         method: "GET".to_string(),
//         url: url.to_string(),
//         headers: HashMap::new(),
//         body: vec![]
//     }
// }
//
// #[test]
// fn test_router() {
//     let mut router = BasicRouter::new();
//     router.add_route("/", |_| get_form_response());
//     router.add_route("/gcd", |req| get_gcd_response(req));
//
//     assert_eq!(router.handle_request(&req("/piano")).code, 404);
//     assert_eq!(router.handle_request(&req("/")).code, 200);
//     assert_eq!(router.handle_request(&req("/gcd")).code, 500);
// }



//                                    copy
// use std::fs;
// use std::io;
// use std::path::Path;
//
// fn copy_dir_to(src: &Path, dst: &Path) -> io::Result<()> {
//     if !dst.is_dir() {
//         fs::create_dir(dst)?;
//     }
//
//     for entry_result in src.read_dir()? {
//         let entry = entry_result?;
//         let file_type = entry.file_type()?;
//         copy_to(&entry.path(), &file_type, &dst.join(entry.file_name()))?;
//     }
//
//     Ok(())
// }
//
// #[cfg(unix)]
// use std::os::unix::fs::symlink;
//
// fn symlink<P: AsRef<Path>, Q: AsRef<Path>>(src: P, _dst: Q) -> std::io::Result<()> {
//     Err(io::Error::new(io::ErrorKind::Other,
//                 format!("cant copy symbolic link: {}", src.as_ref().display())))
// }
//
// fn copy_to(src: &Path, src_type: &fs::FileType, dst: &Path) -> io::Result<()> {
//     if src_type.is_file() {
//         fs::copy(src, dst);
//     } else if src_type.is_dir() {
//         copy_dir_to(src, dst);
//     } else if src_type.is_symlink() {
//         let target = src.read_link()?;
//         symlink(target, dst)?;
//     } else {
//         return Err(io::Error::new(io::ErrorKind::Other,
//                                   format!("don't know how to copy: {}",
//                                           src.display())));
//     }
//
//     Ok(())
// }
//
// fn copy_into<P, Q>(source: P, destination: Q) -> io::Result<()>
//     where P: AsRef<Path>, Q: AsRef<Path>
// {
//     let src = source.as_ref();
//     let dst = destination.as_ref();
//
//     match src.file_name() {
//         None => {
//             return Err(io::Error::new(io::ErrorKind::Other,
//                                         format!("cant copy nameless directory: {}", src.display())));
//         }
//         Some(src_name) => {
//             let md = src.metadata()?;
//             copy_to(src, &md.file_type(), &dst.join(src_name))?;
//         }
//     }
//
//     Ok(())
// }
//
// fn dwim_copy<P, Q>(source: P, destination: Q) -> io::Result<()>
//     where P: AsRef<Path>,
//           Q: AsRef<Path>
// {
//     let src = source.as_ref();
//     let dst = destination.as_ref();
//
//     if dst.is_dir() {
//         copy_into(src, dst)
//     } else {
//         let md = src.metadata()?;
//         copy_to(src, &md.file_type(), dst)
//     }
// }
//
// fn copy_main() -> io::Result<()> {
//     let args = std::env::args_os().collect::<Vec<_>>();
//     if args.len() < 3 {
//         println!("usage: copy FILE... DESTINATION");
//     } else if args.len() == 3 {
//         dwim_copy(&args[1], &args[2])?;
//     } else {
//         let dst = Path::new(&args[args.len()-1]);
//         if !dst.is_dir() {
//             return Err(io::Error::new(io::ErrorKind::Other,
//             format!("target '{}' is not a directory", dst.display())));
//         }
//         for i in 1 .. args.len()-1 {
//             copy_into(&args[i], dst)?;
//         }
//     }
//
//     Ok(())
// }
//
// fn main() {
//     use std::io::Write;
//
//     if let Err(err) = copy_main() {
//         writeln!(io::stderr(), "error: {}", err).unwrap();
//     }
// }



//                                          echo-server
// use std::net::TcpListener;
// use std::io;
// use std::thread::spawn;
//
// fn echo_main(addr: &str) -> io::Result<()> {
//     let listener = TcpListener::bind(addr)?;
//     println!("listening on {}", addr);
//     loop {
//         let (mut stream, addr) = listener.accept()?;
//         println!("connection received from {}", addr);
//
//         let mut write_stream = stream.try_clone()?;
//         spawn(move || {
//             io::copy(&mut stream, &mut write_stream)
//                 .expect("error in client thread: ");
//             println!("connection closed");
//         })
//     }
// }
//
// fn main() {
//     echo_main("127.0.0.1").expect("error: ");
// }




//                                      gcd
// fn gcd(mut n: u64, mut m: u64) -> u64 {
//     assert!(n != 0 && m != 0);
//     while m != 0 {
//         if m < n {
//             let t = m;
//             m = n;
//             n = t;
//         }
//         m = m % n
//     }
//     n
// }
//
// #[test]
// fn test_gcd() {
//     assert_eq!(gcd(14, 15), 1);
//
//     assert_eq!(gcd(2 * 3 * 5 * 11 * 17,
//                    3 * 7 * 11 * 13 * 19),
//                3 * 11);
// }
//
// use std::io::Write;
// use std::str::FromStr;
//
// fn main() {
//     let mut numbers = Vec::new();
//     for arg in std::env::args().skip(1) {
//         numbers.push(u64::from_str(&arg).expect("error parsing arg"));
//     }
//
//     if numbers.len() == 0 {
//         writeln!(std::io::stderr(), "Usage: gcd NUMBER ...").unwrap();
//         std::process::exit(1);
//     }
//
//     let mut d = numbers[0];
//     for m in &numbers[1..] {
//         d = gcd(d, *m);
//     }
//
//     println!("The greatest common divisor of {:?} is {}",
//              numbers, d);
// }



//                      http-get
// use std::error::Error;
// use std::io;
//
// fn http_get_main(url: &str) -> Result<(), Box<dyn Error>> {
//     let mut response = reqwest::get(url)?;
//     if !response.status().is_success() {
//         Err(format!("{}", response.status()))?;
//     }
//
//     let stdout = io::stdout();
//     io::copy(&mut response, &mut stdout.lock())?;
//
//     Ok(())
// }
//
// fn main() {
//     let args: Vec<String> = std::env::args().collect();
//     if args.len() != 2 {
//         eprintln!("usage: http-get URL");
//         return;
//     }
//
//     if let Err(err) = http_get_main(&args[1]) {
//         eprintln!("error: {}", err);
//     }
// }



//              liggit2
extern crate libc;

fn main() {
    let path = std::env::args_os().skip(1).next()
        .expect("usage: libgit2-rs PATH");

    let repo = git::Repository::open(&path)
        .expect("opening repository");

    let commit_oid = repo.reference_name_to_id("HEAD")
        .expect("looking up 'HEAD' reference");

    let commit = repo.find_commit(&commit_oid)
        .expect("looking up commit");

    let author = commit.author();
    println!("{} <{}>\n",
             author.name().unwrap_or("(none)"),
             author.email().unwrap_or("none"));

    println!("{}", commit.message().unwrap_or("(none)"));
}