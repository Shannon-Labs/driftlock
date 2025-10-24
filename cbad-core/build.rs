// Build script for cbad-core
// Handles linking OpenZL C library when openzl feature is enabled

fn main() {
    #[cfg(feature = "openzl")]
    {
        let openzl_dir = std::path::PathBuf::from(env!("CARGO_MANIFEST_DIR"))
            .parent()
            .expect("parent directory exists")
            .join("deps/openzl");

        println!("cargo:rerun-if-changed={}", openzl_dir.display());

        // Link search path for OpenZL library
        println!("cargo:rustc-link-search=native={}", openzl_dir.display());

        // Link OpenZL static library
        println!("cargo:rustc-link-lib=static=openzl");

        // Link OpenZL's bundled zstd
        let zstd_lib = openzl_dir.join("deps/zstd/lib");
        println!("cargo:rustc-link-search=native={}", zstd_lib.display());
        println!("cargo:rustc-link-lib=static=zstd");

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
        println!("cargo:include={}", include_path.display());
    }
}
