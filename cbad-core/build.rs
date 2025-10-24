// Build script for cbad-core
// Handles linking OpenZL C library when openzl feature is enabled

use std::env;
use std::path::PathBuf;

fn main() {
    #[cfg(feature = "openzl")]
    {
        let manifest_dir = PathBuf::from(env::var("CARGO_MANIFEST_DIR").unwrap());
        let openzl_dir = manifest_dir.parent().expect("parent directory exists").join("openzl");
        
        // Find the actual static library location
        let lib_dir = if openzl_dir.join("libopenzl.a").exists() {
            openzl_dir.clone()
        } else {
            // Check cached build directories for the static library
            let cached_objs = openzl_dir.join("cachedObjs");
            if cached_objs.exists() {
                let mut found_lib = None;
                for entry in std::fs::read_dir(&cached_objs).unwrap() {
                    let entry = entry.unwrap();
                    let path = entry.path();
                    if path.is_dir() && path.join("libopenzl.a").exists() {
                        found_lib = Some(path);
                        break;
                    }
                }
                found_lib.unwrap_or_else(|| {
                    eprintln!("Warning: Could not find libopenzl.a, using main openzl directory");
                    openzl_dir.clone()
                })
            } else {
                openzl_dir.clone()
            }
        };

        println!("cargo:rerun-if-changed={}", lib_dir.display());
        
        // Print the library location for debugging
        println!("cargo:warning=OpenZL static library found at: {}", lib_dir.display());

        // Link search path for OpenZL library
        println!("cargo:rustc-link-search=native={}", lib_dir.display());

        // Link OpenZL static library explicitly - use absolute path to avoid dynamic linking
        println!("cargo:rustc-link-lib=static=openzl");

        // Link zstd library - try system first, then bundled
        let system_zstd = PathBuf::from("/opt/homebrew/Cellar/zstd/1.5.7/lib");
        let bundled_zstd = openzl_dir.join("deps/zstd/lib");
        
        if system_zstd.exists() {
            println!("cargo:rustc-link-search=native={}", system_zstd.display());
            println!("cargo:rustc-link-lib=dylib=zstd");
            println!("cargo:rustc-link-arg=-Wl,-rpath,{}", system_zstd.display());
        } else if bundled_zstd.exists() {
            println!("cargo:rustc-link-search=native={}", bundled_zstd.display());
            println!("cargo:rustc-link-lib=static=zstd");
            println!("cargo:rustc-link-arg=-Wl,-rpath,{}", bundled_zstd.display());
        }

        // Link C++ standard library (OpenZL uses C++)
        // Use different lib name depending on platform
        #[cfg(target_os = "linux")]
        println!("cargo:rustc-link-lib=dylib=stdc++");

        #[cfg(target_os = "macos")]
        println!("cargo:rustc-link-lib=dylib=c++");

        // Link pthread (required by OpenZL)
        println!("cargo:rustc-link-lib=dylib=pthread");

        // Include path for OpenZL headers
        let include_path = openzl_dir.join("include");
        if include_path.exists() {
            println!("cargo:include={}", include_path.display());
        }
    }
}
