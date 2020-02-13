extern crate clap;
use clap::App;

mod create;
mod init;

fn main() {
    let matches = App::new("Jot")
        .version("0.2")
        .author("Joel Holmes <holmes89@gmail.com>")
        .about("Create quick notes")
        .subcommand(
            App::new("create").about("create new entry"), //TODO add functionality for "collections"
        )
        .subcommand(App::new("edit").about("edit existing entry")) //TODO default to today, otherwise allow for edit entry date or collection
        .subcommand(App::new("read").about("read entry")) //TODO default to today, otherwise allow for read entry date or collection
        .subcommand(App::new("init").about("initialize app")) // initialize application
        .get_matches();

    match matches.subcommand_name() {
        Some("create") => crate::create::new_entry(),
        Some("edit") => println!("edited"),
        Some("read") => println!("read"),
        Some("init") => crate::init::intialize(),
        None => println!("HALP!"),
        _ => unreachable!(),
    }
}
