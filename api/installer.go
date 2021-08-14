package api

// get stats of the APP
// control STOP
// releases /releases/<string:service>
// RETURN
//[
//    "v1.7.2",
//    "v1.7.1"
//]


/*

   def get_releases_link(self) -> str:
	  return 'https://api.github.com/repos/NubeIO/{}/releases'.format(self.repo_name)

   def get_download_link(self, token: str, is_browser_downloadable: bool = False):
       headers = {}
       if token:
           headers['Authorization'] = f'Bearer {token}'
       release_link: str = f'https://api.github.com/repos/NubeIO/{self.repo_name}/releases/tags/{self.version}'
       resp = requests.get(release_link, headers=headers)
       row: str = json.loads(resp.content)
       setting: AppSetting = current_app.config[AppSetting.FLASK_KEY]
       download_link = self.select_link(row, is_browser_downloadable)
       if not download_link:
           raise ModuleNotFoundError(
               f'No app {self.service} for type {setting.device_type} & version {self.version}, '
               f'check your token & repo')
       return download_link

   def get_latest_release(self, token: str):
       headers = {}
       if token:
           headers['Authorization'] = f'Bearer {token}'
       release_link: str = f'https://api.github.com/repos/NubeIO/{self.repo_name}/releases'
       resp = requests.get(release_link, headers=headers)
       data = json.loads(resp.content)
       latest_release = ''
       for row in data:
           if isinstance(row, str):
               raise PreConditionException('Please insert GitHub valid token!')
           release = row.get('tag_name')
           if not latest_release or packaging_version.parse(latest_release) <= packaging_version.parse(release):
               latest_release = release
       if not latest_release:
           raise NotFoundException('Latest release not found!')
       return latest_release

 */
