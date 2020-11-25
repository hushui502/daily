use rustproject::{NewArticle, Summary, Tweet};

fn main() {
    let article = NewArticle {
        headline: String::from("new headline"),
        author: String::from("hufan"),
        location: String::from("china"),
        content: "".to_string(),
    };

    println!("1 new article: {}", article.summarize());

    notify(article)

}

pub fn notify(item: impl Summary) {
    println!("breaking news! {}", item.summarize())
}