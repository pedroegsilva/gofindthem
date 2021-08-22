load("@bazel_gazelle//:deps.bzl", "go_repository")

def go_dependencies():
    go_repository(
        name = "com_github_anknown_ahocorasick",
        importpath = "github.com/anknown/ahocorasick",
        sum = "h1:onfun1RA+KcxaMk1lfrRnwCd1UUuOjJM/lri5eM1qMs=",
        version = "v0.0.0-20190904063843-d75dbd5169c0",
    )
    go_repository(
        name = "com_github_anknown_darts",
        importpath = "github.com/anknown/darts",
        sum = "h1:HblK3eJHq54yET63qPCTJnks3loDse5xRmmqHgHzwoI=",
        version = "v0.0.0-20151216065714-83ff685239e6",
    )
    go_repository(
        name = "com_github_cloudflare_ahocorasick",
        importpath = "github.com/cloudflare/ahocorasick",
        sum = "h1:8yL+85JpbwrIc6m+7N1iYrjn/22z68jwrTIBOJHNe4k=",
        version = "v0.0.0-20210425175752-730270c3e184",
    )
    go_repository(
        name = "com_github_davecgh_go_spew",
        importpath = "github.com/davecgh/go-spew",
        sum = "h1:ZDRjVQ15GmhC3fiQ8ni8+OwkZQO4DARzQgrnXU1Liz8=",
        version = "v1.1.0",
    )
    go_repository(
        name = "com_github_petar_dambovaliev_aho_corasick",
        importpath = "github.com/petar-dambovaliev/aho-corasick",
        sum = "h1:WuXe30Ig5zUIYEHyzsLMBFPP5l0yRQ5IiZScODHwy8g=",
        version = "v0.0.0-20210512121028-af76a9ff7276",
    )
    go_repository(
        name = "com_github_pmezard_go_difflib",
        importpath = "github.com/pmezard/go-difflib",
        sum = "h1:4DBwDE0NGyQoBHbLQYPwSUPoCMWR5BEzIk/f1lZbAQM=",
        version = "v1.0.0",
    )
    go_repository(
        name = "com_github_stretchr_objx",
        importpath = "github.com/stretchr/objx",
        sum = "h1:4G4v2dO3VZwixGIRoQ5Lfboy6nUhCyYzaqnIAPPhYs4=",
        version = "v0.1.0",
    )
    go_repository(
        name = "com_github_stretchr_testify",
        importpath = "github.com/stretchr/testify",
        sum = "h1:nwc3DEeHmmLAfoZucVR881uASk0Mfjw8xYJ99tb5CcY=",
        version = "v1.7.0",
    )
    go_repository(
        name = "in_gopkg_check_v1",
        importpath = "gopkg.in/check.v1",
        sum = "h1:yhCVgyC4o1eVCa2tZl7eS0r+SDo693bJlVdllGtEeKM=",
        version = "v0.0.0-20161208181325-20d25e280405",
    )
    go_repository(
        name = "in_gopkg_yaml_v3",
        importpath = "gopkg.in/yaml.v3",
        sum = "h1:dUUwHk2QECo/6vqA44rthZ8ie2QXMNeKRTHCNY2nXvo=",
        version = "v3.0.0-20200313102051-9f266ea9e77c",
    )
