class RbacWizard < Formula
  desc "RBAC Wizard is a tool that helps you visualize and analyze the RBAC configurations of your Kubernetes cluster. It provides a graphical representation of the Kubernetes RBAC objects."
  homepage "https://github.com/pehlicd/rbac-wizard"
  url "https://github.com/pehlicd/rbac-wizard.git",
      tag:      "v0.0.4",
      revision: "827ba7f5ecfaabaf02481d59e198651ad7be259c"
  license "MIT"
  head "https://github.com/pehlicd/rbac-wizard.git", branch: "main"

  depends_on "go" => :build

  def install
    project = "github.com/pehlicd/rbac-wizard"
    ldflags = %W[
      -s -w
      -X #{project}/cmd.versionString=#{version}
      -X #{project}/cmd.buildCommit=#{Utils.git_head}
      -X #{project}/cmd.buildDate=#{time.iso8601}
    ]
    system "go", "build", *std_go_args(ldflags: ldflags)
  end

  test do
    assert_match version.to_s, "#{bin}/rbac-wizard version"
  end
end