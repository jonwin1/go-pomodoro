{
    description = "A simple pomodoro program intended for use with Waybar";

    inputs.nixpkgs.url = "nixpkgs/nixos-24.11";

    outputs = { self, nixpkgs }: let
        forAllSystems = function:
            nixpkgs.lib.genAttrs [
                "x86_64-linux"
                "x86_64-darwin"
                "aarch64-linux"
                "aarch64-darwin"
            ] (system: function nixpkgs.legacyPackages.${system});
    in {
        packages = forAllSystems (pkgs: {
            default = pkgs.buildGoModule {
                pname = "pomodoro";
                version = self.lastModifiedDate;
                src = ./.;
                nativeBuildInputs = with pkgs; [ pkg-config ];
                buildInputs = with pkgs; [ alsa-lib ];
                # vendorHash = pkgs.lib.fakeHash;
                vendorHash = "sha256-z/ZM97Ti9hYF+B8EbwvqXKlgzzd2MBBeiNtA9wH1jwU=";
            };
        });

        devShells = forAllSystems (pkgs: {
            default = pkgs.mkShell {
                buildInputs = with pkgs; [
                    go
                    gopls
                    go-tools
                    gotools

                    pkg-config
                    alsa-lib
                ];
            };
        });
    };
}
