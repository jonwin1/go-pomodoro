{
  description = "A simple pomodoro program intended for use with Waybar";

  inputs.nixpkgs.url = "nixpkgs/nixos-24.11";

  outputs =
    { self, nixpkgs }:
    let
      forAllSystems =
        function:
        nixpkgs.lib.genAttrs [
          "x86_64-linux"
          "x86_64-darwin"
          "aarch64-linux"
          "aarch64-darwin"
        ] (system: function nixpkgs.legacyPackages.${system});

      commonDeps = pkgs: {
        # Run-time dependencies
        buildInputs = with pkgs; [
          alsa-lib
        ];

        # Build-time dependencies
        nativeBuildInputs = with pkgs; [
          go
          gopls
          pkg-config
        ];
      };
    in
    {
      packages = forAllSystems (
        pkgs:
        let
          deps = commonDeps pkgs;
        in
        {
          default = pkgs.buildGoModule {
            pname = "pomodoro";
            version = self.shortRev or self.dirtyShortRev;
            src = ./.;

            buildInputs = deps.buildInputs;
            nativeBuildInputs = deps.nativeBuildInputs;

            vendorHash = "sha256-z/ZM97Ti9hYF+B8EbwvqXKlgzzd2MBBeiNtA9wH1jwU=";
          };
        }
      );

      devShells = forAllSystems (
        pkgs:
        let
          deps = commonDeps pkgs;
        in
        {
          default = pkgs.mkShell {
            buildInputs = deps.buildInputs;
            nativeBuildInputs = deps.nativeBuildInputs;
          };
        }
      );
    };
}
