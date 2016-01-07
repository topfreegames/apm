# config valid only for current version of Capistrano
lock '3.4.0'

set :tfg_go_repo, "github.com/topfreegames"
set :apm_go_repo, "github.com/topfreegames/apm"
set :application, 'apm'
set :repo_url, 'https://github.com/topfreegames/apm.git'

set :stages, %w(staging production)
set :default_stage, "staging"

set :keep_releases, 20
set :scm, :git
set :format, :pretty
set :pty, true

set :deploy_to, "/var/apps/#{fetch :application}/#{fetch :stage}"
set :shared_path, "#{fetch :deploy_to}/shared"

set :gopath, "#{fetch :shared_path}/gopath"
set :gobin, "/usr/local/go/bin"
set :goenv, "GOPATH=#{fetch :gopath} GOBIN=#{fetch :gobin} PATH=\"/usr/local/go/bin:$PATH\""

set :go, "#{fetch :goenv} go"

set :goget, "#{fetch :go} get -u -f all || true && #{fetch :go} get"
set :gobuild, "#{fetch :go} build"
set :goinstall, "#{fetch :go} install"

set :apm_config_path, "#{fetch :shared_path}/apm-config/config.toml"

namespace :deploy do

  desc 'Link repo'
  task :link_repo do
    on roles(:app) do
      execute "mkdir -p \"$(dirname \"#{fetch :apm_config_path}\")\" && touch \"#{fetch :apm_config_path}\""
      execute "mkdir -p #{fetch :gopath}/src/#{fetch :tfg_go_repo}"
      execute "ln -snf #{fetch :release_path} #{fetch :gopath}/src/#{fetch :apm_go_repo}"      
    end    
  end
  desc 'Compile'
  task :compile do
    on roles(:app) do
      execute "cd #{fetch :gopath}/src/#{fetch :apm_go_repo} && #{fetch :goget} && #{fetch :gobuild}"
    end
  end
  desc 'Start'
  task :start do
    on roles(:app) do
      execute "cd #{fetch :gopath}/src/#{fetch :apm_go_repo} && ./apm serve --config-file=\"#{fetch :apm_config_path}\""
    end
  end
  
  after :updating, :link_repo
  after :publishing, :compile
  after :compile, :start
end
