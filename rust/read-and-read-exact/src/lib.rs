#[cfg(test)]
mod tests {
    use std::io::{self, prelude::*};

    #[test]
    fn just_size() -> io::Result<()> {
        let mut buf = [0 as u8; 5];
        let size = "abcde".as_bytes().read(&mut buf)?;
        assert_eq!(size, 5);
        assert_eq!(buf, "abcde".as_bytes());

        let mut buf = [0 as u8; 5];
        "abcde".as_bytes().read_exact(&mut buf)?;
        assert_eq!(buf, "abcde".as_bytes());

        Ok(())
    }

    #[test]
    fn small_buf() -> io::Result<()> {
        let src = "abcde";

        let mut buf = [0 as u8; 3];
        let size = src.as_bytes().read(&mut buf)?;
        assert_eq!(size, 3);
        assert_eq!(buf, "abc".as_bytes());

        let mut buf = [0 as u8; 3];
        src.as_bytes().read_exact(&mut buf)?;
        assert_eq!(buf, "abc".as_bytes());

        Ok(())
    }

    #[test]
    fn large_buf() -> io::Result<()> {
        let src = "abcde";

        let mut buf = [0 as u8; 8];
        let size = src.as_bytes().read(&mut buf)?;
        assert_eq!(size, 5);
        assert_eq!(buf, "abcde\0\0\0".as_bytes());

        let mut buf = [0 as u8; 8];
        match src.as_bytes().read_exact(&mut buf) {
            Ok(x) => assert!(false, "got {:?}", x),
            Err(e) => assert_eq!(e.kind(), io::ErrorKind::UnexpectedEof),
        }

        Ok(())
    }
}
