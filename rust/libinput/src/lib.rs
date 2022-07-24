use std::fs::{File, OpenOptions};
use std::os::unix::prelude::*;
use std::path::Path;

use input::LibinputInterface;
use libc::{O_RDONLY, O_RDWR, O_WRONLY};
use nix::ioctl_write_int;

pub struct Interface;

impl LibinputInterface for Interface {
    fn open_restricted(&mut self, path: &Path, flags: i32) -> Result<RawFd, i32> {
        OpenOptions::new()
            .custom_flags(flags)
            .read((flags & O_RDONLY != 0) | (flags & O_RDWR != 0))
            .write((flags & O_WRONLY != 0) | (flags & O_RDWR != 0))
            .open(path)
            .map(|file| file.into_raw_fd())
            .map_err(|err| err.raw_os_error().unwrap())
    }
    fn close_restricted(&mut self, fd: RawFd) {
        unsafe {
            File::from_raw_fd(fd);
        }
    }
}

ioctl_write_int!(eviocgrab, b'E', 0x90);

pub fn grab<T: AsRawFd>(fd: T) -> Result<(), i32> {
    let result: nix::Result<libc::c_int> = unsafe { eviocgrab(fd.as_raw_fd(), 1) };
    result.map_err(|e| e as i32).map(|_| ())
}
