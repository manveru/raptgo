task :default do
  files = Dir['*.go']
  imports = %w[
    /home/manveru/github/banthar/Go-SDL
    /home/manveru/github/banthar/Go-SDL/sdl
  ].map{|i| ['-I', i] }.flatten
  sh("6g", "-o", "raptgo.6", *imports, *files)
  sh("6l -o raptgo raptgo.6")
end
