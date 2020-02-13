extern crate dirs;
extern crate rustyline;

use rustyline::error::ReadlineError;
use rustyline::Editor;
use std::path::Path;

pub fn intialize() {
    let mut rl = Editor::<()>::new();
    let readline = rl.readline("Remote Git URL >> ");
    match readline {
        Ok(url) => println!("{}", clone_repo(&url)),
        Err(ReadlineError::Interrupted) => {
            println!("CTRL-C");
        }
        Err(ReadlineError::Eof) => {
            println!("CTRL-D");
        }
        Err(err) => {
            println!("Error: {:?}", err);
        }
    }
}

fn clone_repo(url: &str) -> String {
    let home_dir = dirs::home_dir().expect("cannot find home directory");
    let home = home_dir.to_str().expect("could not extract directory name");
    let jot_home = format!("{}{}", home, "/.jot");
    let _pub_path = format!("{}{}", home, "/.ssh/id_rsa.pub");
    let _private_path = format!("{}{}", home, "/.ssh/id_rsa");
    let p = Path::new(&jot_home);
    if p.exists() {
        panic!("directory already exists")
    }

    let mut repo_builder = git2::build::RepoBuilder::new();
    let mut fetch_options = git2::FetchOptions::new();
    let mut auth_callback = git2::RemoteCallbacks::new();

    auth_callback.credentials(|_, _, _| {
        let credentials =
            git2::Cred::ssh_key_from_agent("git").expect("Could not create credentials object");
        Ok(credentials)
    });

    fetch_options.remote_callbacks(auth_callback);
    repo_builder.fetch_options(fetch_options);
    let _repo = repo_builder.clone(url, p).expect("could not clone repo");
    String::from("initialized")
}
