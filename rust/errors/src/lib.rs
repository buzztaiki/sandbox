pub mod use_option {
    use std::{fs, io::Read};

    pub fn with_match(fname: &str) -> Option<i64> {
        match fs::File::open(fname) {
            Ok(mut f) => {
                let mut buf = String::new();
                match f.read_to_string(&mut buf) {
                    Ok(_) => match buf.trim().parse::<i64>() {
                        Ok(n) => Some(n),
                        Err(e) => {
                            eprintln!("failed to parse: {}", e);
                            None
                        }
                    },
                    Err(e) => {
                        eprintln!("failed to read: {}", e);
                        None
                    }
                }
            }
            Err(e) => {
                eprintln!("failed to open: {}", e);
                None
            }
        }
    }

    pub fn with_early_return(fname: &str) -> Option<i64> {
        let mut f = match fs::File::open(fname) {
            Ok(f) => f,
            Err(e) => {
                eprintln!("failed to open: {}", e);
                return None;
            }
        };

        let mut buf = String::new();
        match f.read_to_string(&mut buf) {
            Ok(_) => {}
            Err(e) => {
                eprintln!("failed to read: {}", e);
                return None;
            }
        };

        let n = match buf.trim().parse::<i64>() {
            Ok(n) => n,
            Err(e) => {
                eprintln!("failed to parse: {}", e);
                return None;
            }
        };

        Some(n)
    }

    pub fn with_result_map_or(fname: &str) -> Option<i64> {
        let mut f = fs::File::open(fname).map_or(None, |f| Some(f))?;
        let mut buf = String::new();
        let s = f.read_to_string(&mut buf).map_or(None, |_| Some(buf))?;
        s.trim().parse::<i64>().map_or(None, |n| Some(n))
    }

    pub fn with_result_ok(fname: &str) -> Option<i64> {
        let mut f = fs::File::open(fname).ok()?;
        let mut buf = String::new();
        f.read_to_string(&mut buf).ok()?;
        let s = buf;
        s.trim().parse::<i64>().ok()
    }

    #[cfg(test)]
    mod tests {
        macro_rules! assert_use_option {
            ($f:expr) => {
                assert_eq!($f("/dev/null"), None);
                assert_eq!($f("testfiles/string.txt"), None);
                assert_eq!($f("testfiles/number.txt"), Some(42));
            };
        }

        #[test]
        fn with_match() {
            assert_use_option!(super::with_match);
        }

        #[test]
        fn with_early_return() {
            assert_use_option!(super::with_early_return);
        }

        #[test]
        fn with_result_map_or() {
            assert_use_option!(super::with_result_map_or);
        }

        #[test]
        fn with_result_ok() {
            assert_use_option!(super::with_result_ok);
        }
    }
}

pub mod use_result {
    use std::{error, fmt, fs::File, io, io::Read, num};
    use anyhow::Context;

    pub fn with_box_error(fname: &str) -> Result<i64, Box<dyn error::Error>> {
        let mut f = File::open(fname)?;
        let mut buf = String::new();
        f.read_to_string(&mut buf)?;
        let s = buf;
        let n = s.trim().parse::<i64>()?;
        Ok(n)
    }

    #[derive(Debug)]
    pub enum Error {
        IoError(io::Error),
        ParseIntError(num::ParseIntError),
    }

    impl error::Error for Error {}

    impl fmt::Display for Error {
        fn fmt(&self, f: &mut fmt::Formatter<'_>) -> fmt::Result {
            match self {
                Error::IoError(e) => e.fmt(f),
                Error::ParseIntError(e) => e.fmt(f),
            }
        }
    }

    impl From<io::Error> for Error {
        fn from(e: io::Error) -> Self {
            Error::IoError(e)
        }
    }

    impl From<num::ParseIntError> for Error {
        fn from(e: num::ParseIntError) -> Self {
            Error::ParseIntError(e)
        }
    }

    pub fn with_enum_error(fname: &str) -> Result<i64, Error> {
        let mut f = File::open(fname)?;
        let mut buf = String::new();
        f.read_to_string(&mut buf)?;
        let s = buf;
        let n = s.trim().parse::<i64>()?;
        Ok(n)
    }

    pub fn with_anyhow(fname: &str) -> anyhow::Result<i64> {
        let mut f = File::open(fname)?;
        let mut buf = String::new();
        f.read_to_string(&mut buf)?;
        let s = buf;
        let n = s.trim().parse::<i64>().with_context(|| format!("failed to parse '{}'", s.trim()))?;
        Ok(n)
    }

    #[derive(thiserror::Error, Debug)]
    pub enum Error2 {
        #[error("{0}")]
        IoError(#[from] io::Error),
        #[error("{0}")]
        ParseIntError(#[from] num::ParseIntError),
    }

    pub fn with_thiserror(fname: &str) -> Result<i64, Error2> {
        let mut f = File::open(fname)?;
        let mut buf = String::new();
        f.read_to_string(&mut buf)?;
        let s = buf;
        let n = s.trim().parse::<i64>()?;
        Ok(n)
    }

    #[cfg(test)]
    mod tests {
        use std::io;

        #[test]
        fn with_box_error() {
            match super::with_box_error("/dev/unknown") {
                Err(e) => assert_eq!(e.to_string(), "No such file or directory (os error 2)"),
                x => assert!(false, "got {:?}", x),
            };
            match super::with_box_error("/dev/null") {
                Err(e) => assert_eq!(e.to_string(), "cannot parse integer from empty string"),
                x => assert!(false, "got {:?}", x),
            };
            match super::with_box_error("testfiles/string.txt") {
                Err(e) => assert_eq!(e.to_string(), "invalid digit found in string"),
                x => assert!(false, "got {:?}", x),
            };

            match super::with_box_error("testfiles/number.txt") {
                Ok(n) => assert_eq!(n, 42),
                x => assert!(false, "got {:?}", x),
            };

            // down cast で元のエラーは取得できるがちと面倒
            match super::with_box_error("/dev/unknown") {
                Err(e) => {
                    match e.downcast::<io::Error>() {
                        Ok(e) => assert_eq!(e.to_string(), "No such file or directory (os error 2)"),
                        x => assert!(false, "got {:?}", x),
                    }
                }
                x => assert!(false, "got {:?}", x),
            };
        }

        #[test]
        fn with_enum_error() {
            match super::with_enum_error("/dev/unknown") {
                Err(super::Error::IoError(e)) => assert_eq!(e.kind(), io::ErrorKind::NotFound),
                x => assert!(false, "got {:?}", x),
            };
            match super::with_enum_error("/dev/null") {
                Err(super::Error::ParseIntError(e)) => assert_eq!(e.to_string(), "cannot parse integer from empty string"),
                x => assert!(false, "got {:?}", x),
            };
            match super::with_enum_error("testfiles/string.txt") {
                Err(super::Error::ParseIntError(e)) => assert_eq!(e.to_string(), "invalid digit found in string"),
                x => assert!(false, "got {:?}", x),
            };
            match super::with_enum_error("testfiles/number.txt") {
                Ok(n) => assert_eq!(n, 42),
                x => assert!(false, "got {:?}", x),
            };
        }

        #[test]
        fn with_anyhow() {
            match super::with_anyhow("/dev/unknown") {
                Err(e) => assert_eq!(e.to_string(), "No such file or directory (os error 2)"),
                x => assert!(false, "got {:?}", x),
            };
            match super::with_anyhow("/dev/null") {
                Err(e) => assert_eq!(e.to_string(), "failed to parse ''"),
                x => assert!(false, "got {:?}", x),
            };
            match super::with_anyhow("testfiles/string.txt") {
                Err(e) => assert_eq!(e.to_string(), "failed to parse 'moo'"),
                x => assert!(false, "got {:?}", x),
            };

            match super::with_anyhow("testfiles/number.txt") {
                Ok(n) => assert_eq!(n, 42),
                x => assert!(false, "got {:?}", x),
            };

            // down cast で元のエラーは取得できるがちと面倒
            match super::with_anyhow("/dev/unknown") {
                Err(e) => {
                    match e.downcast::<io::Error>() {
                        Ok(e) => assert_eq!(e.to_string(), "No such file or directory (os error 2)"),
                        x => assert!(false, "got {:?}", x),
                    }
                }
                x => assert!(false, "got {:?}", x),
            };
        }

        #[test]
        fn with_thiserror() {
            match super::with_thiserror("/dev/unknown") {
                Err(super::Error2::IoError(e)) => assert_eq!(e.kind(), io::ErrorKind::NotFound),
                x => assert!(false, "got {:?}", x),
            };
            match super::with_thiserror("/dev/null") {
                Err(super::Error2::ParseIntError(e)) => assert_eq!(e.to_string(), "cannot parse integer from empty string"),
                x => assert!(false, "got {:?}", x),
            };
            match super::with_thiserror("testfiles/string.txt") {
                Err(super::Error2::ParseIntError(e)) => assert_eq!(e.to_string(), "invalid digit found in string"),
                x => assert!(false, "got {:?}", x),
            };
            match super::with_thiserror("testfiles/number.txt") {
                Ok(n) => assert_eq!(n, 42),
                x => assert!(false, "got {:?}", x),
            };
        }

    }
}
