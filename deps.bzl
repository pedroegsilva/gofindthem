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
        name = "com_github_petar_dambovaliev_aho_corasick",
        importpath = "github.com/petar-dambovaliev/aho-corasick",
        sum = "h1:WuXe30Ig5zUIYEHyzsLMBFPP5l0yRQ5IiZScODHwy8g=",
        version = "v0.0.0-20210512121028-af76a9ff7276",
    )
