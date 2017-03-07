#!/usr/bin/env ruby

require "yaml"
require "fileutils"

class Project
  include FileUtils
  attr_accessor :url
  def initialize(url)
    @url = url
  end

  def project_name
    UrlExtractor.new(url).run
  end

  def full_path
    File.join("src", project_name)
  end

  def exists?
    File.exists?(full_path)
  end

  def refresh!
    cd(full_path) do
      exec("git pull")
    end
  end

  def checkout!
    exec("git clone #{url} #{full_path}")
  end

  def clear!
    if exists?
      rm_rf(full_path)
    end
  end

  private

  def exec(cmd)
    puts system(cmd)
  end
end

class UrlExtractor
  def initialize(url); @url = url ; end
  def run
    parts = clean_url.split("/")
    user  = parts[-2]
    repo  = parts[-1].gsub(/\.git$/, "")
    "#{user}--#{repo}"
  end

  def clean_url
    @url.gsub(":", "/")
  end
end

class ProjectRefresher
  attr_accessor :project
  def initialize(url)
    @project = Project.new(url)
  end

  def run
    if project.exists?
      puts "refresh #{project.full_path}"
      project.refresh!
    else
      puts "checkout #{project.full_path}"
      project.checkout!
    end
  end

end

class Downloader
  def run
    urls.map{|u| for_project(u) }
  end

  def for_project(url)
    ProjectRefresher.new(url).run
  end

  def urls
    @urls ||= YAML.load_file(config_file)["urls"]
  end

  def config_file
    File.join(__dir__, "urls.yml")
  end
end

Downloader.new.run
