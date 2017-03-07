#!/usr/bin/env ruby
require "fileutils"
urls = []
Dir.glob("src/*").each do |f|
  FileUtils.cd(f) do
    content = File.read(".git/config")
    url = content.split("\n").grep(/url/).first
    urls << url.split("= ").last
  end
end

puts urls.map{|x| "- #{x}"}
