class RbacWizard < Formula
  desc "RBAC Wizard is a tool that helps you visualize and analyze the RBAC configurations of your Kubernetes cluster. It provides a graphical representation of the Kubernetes RBAC objects."
  homepage "https://github.com/pehlicd/rbac-wizard"
  url "https://github.com/pehlicd/rbac-wizard.git",
      tag:      "v0.0.3",
      revision: "9b4041cd7d83778ffbda8eb6f5314e5b0f449f9d"
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