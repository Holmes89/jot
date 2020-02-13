use chrono::prelude::*;
use rand::Rng;
use std::{fs::remove_file, fs::File, io::Read, io::Write, process::Command};

pub fn new_entry() {
    let editor = "emacs";
    let mut rng = rand::thread_rng();
    let fname = format!("/tmp/jot-entry-{}", rng.gen_range(0, 1000));

    Command::new(editor)
        .arg(fname.clone())
        .status()
        .expect("unable to create temporary file");

    let mut editable = String::new();
    File::open(fname.clone())
        .expect("Could not open file")
        .read_to_string(&mut editable)
        .expect("unable to read file");

    remove_file(fname.clone()).expect("unable to remove temporary file");

    let home_dir = dirs::home_dir().expect("cannot find home directory");
    let home = home_dir.to_str().expect("could not extract directory name");
    let jot_home = format!("{}{}", home, "/.jot");

    let dt = Utc::now();
    let entry_name = format!("{}/entries/{}.md", jot_home, dt.format("%Y-%m-%d"));
    let mut entry_file = File::create(entry_name.clone()).expect("Could not open file");
    write!(entry_file, "{}", editable.as_str()).expect("unable to write entry");
}
