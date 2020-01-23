extern crate clap;
use clap::{App, Arg};

fn main() {
    let _matches = App::new("Jot")
        .version("0.2")
        .author("Joel Holmes <holmes89@gmail.com>")
        .about("Create quick notes")
        .arg(Arg::with_name("create").short("c").long("create"))
        .get_matches();
}
