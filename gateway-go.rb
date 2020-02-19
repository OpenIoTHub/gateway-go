class GatewayGo < Formula
  desc "OpenIoTHub GateWay Client"
  homepage "https://github.com/OpenIoTHub/gateway-go"
  url "https://github.com/OpenIoTHub/gateway-go.git",
      :tag      => "v0.1.38",
      :revision => "ccb16d41730884a3de5caf7e5e8f2ef2157ce596"

  depends_on "go" => :build

  def install
    system "go", "build", "-ldflags",
             "-s -w -X main.version=#{version} -X main.commit=#{stable.specs[:revision]} -X main.builtBy=homebrew",
             "-o", bin/"gateway-go"
  end

  test do
    system "#{bin}/gateway-go", "-v"
  end
end
