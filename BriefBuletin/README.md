# BriefBuletin — Angular frontend

BriefBuletin is an Angular frontend for the BriefBuletin API. It provides a UI for managing articles, categories and user authentication.

**This README focuses on local development, hiding the backend host behind a proxy, and deploying to Vercel while keeping the backend host private.**

## Quick start (development)

**Requirements**:
- Node.js (16+ recommended)
- npm
- Angular CLI (optional for global `ng`)

1. Install dependencies:

```pwsh
npm install
```

2. Start the dev server (recommended: use a proxy so the real backend host is not exposed):

```pwsh
# If you have a proxy configured (see below)
npm start
# otherwise fallback to ng serve
ng serve
```

Open `http://localhost:4200`.

## Hiding the backend host (recommended)

There are two common approaches to avoid exposing the backend host/IP to the browser:

- Local dev proxy: the Angular dev server forwards `/api/*` to the real backend. The browser only sees `localhost:4200/api/...`.
- Production hosting rewrite/proxy: your hosting provider (Vercel in this repo) rewrites `/api/*` server-side to the backend so the browser never sees the real host.

Suggested configuration (examples):

- `src/environments/environment.ts` and `src/environments/environment.prod.ts` should expose an `apiUrl` like:

```ts
export const environment = {
	production: false,
	apiUrl: '/api/'
}
```

- `proxy.conf.json` (local development) example:

```json
{
	"/api": {
		"target": "http://118.67.213.45:8339",
		"secure": false,
		"changeOrigin": true,
		"logLevel": "info"
	}
}
```

Add `--proxy-config proxy.conf.json` to `ng serve` (e.g. in `package.json` `start` script):

```json
"start": "ng serve --proxy-config proxy.conf.json"
```

Now `http://localhost:4200/api/...` will be forwarded to `http://118.67.213.45:8339/...` by the dev server.

## Deploying to Vercel (keep backend hidden)

Vercel can rewrite or proxy routes server-side. Example `vercel.json` rewrite:

```json
{
	"rewrites": [
		{ "source": "/api/(.*)", "destination": "http://118.67.213.45:8339/$1" }
	]
}
```

When deployed, the browser will request `/api/...` from the Vercel-hosted frontend; Vercel will internally call the backend and return the response. The client never sees the backend IP.

Important: your backend is HTTP. Browsers block mixed content when a page served over HTTPS (Vercel) tries to call HTTP directly. Using a server-side proxy/rewriter (like Vercel rewrites) avoids that because the request to the backend originates from Vercel (server-side).

## Build & deploy

Build locally:

```pwsh
npm run build -- --configuration production
# or
ng build --configuration production
```

Deploy to Vercel using the Vercel CLI or the Git integration.

## Testing

Run unit tests:

```pwsh
npm test
```

## Troubleshooting

- If API calls return CORS errors during local development, ensure the dev proxy is configured and used (`ng serve --proxy-config proxy.conf.json`).
- If production requests fail after deploying to Vercel, check the Vercel deployment logs and verify the `vercel.json` rewrite is present in the published commit or project settings.

## Security recommendation

- Prefer serving your backend over HTTPS. If possible, enable HTTPS on the backend and point production rewrites to `https://...` to simplify security and reduce risk.

## Where to look in this repo

- `src/environments/environment.ts` and `environment.prod.ts` — frontend base API URL
- `proxy.conf.json` — (optional) dev proxy configuration
- `vercel.json` — (optional) Vercel rewrite configuration
- `package.json` — dev and build scripts

## Next steps I can help with

- Update `environment.*.ts` and add `proxy.conf.json` and `vercel.json` for you.
- Scan the codebase for hardcoded backend URLs and replace them with `environment.apiUrl`.
- Prepare a Vercel preview deployment and help debug any CORS/mixed-content issues.

---

If you want, I can now add the proxy and Vercel files and update `environment.*.ts` so the app uses `/api/` consistently. Which would you like me to do next?
# stock-management-ui



## Getting started

To make it easy for you to get started with GitLab, here's a list of recommended next steps.

Already a pro? Just edit this README.md and make it your own. Want to make it easy? [Use the template at the bottom](#editing-this-readme)!

## Add your files

- [ ] [Create](https://docs.gitlab.com/ee/user/project/repository/web_editor.html#create-a-file) or [upload](https://docs.gitlab.com/ee/user/project/repository/web_editor.html#upload-a-file) files
- [ ] [Add files using the command line](https://docs.gitlab.com/ee/gitlab-basics/add-file.html#add-a-file-using-the-command-line) or push an existing Git repository with the following command:

```
cd existing_repo
git remote add origin https://gitlab.com/satcominternal/stock-management-ui.git
git branch -M main
git push -uf origin main
```

## Integrate with your tools

- [ ] [Set up project integrations](https://gitlab.com/satcominternal/stock-management-ui/-/settings/integrations)

## Collaborate with your team

- [ ] [Invite team members and collaborators](https://docs.gitlab.com/ee/user/project/members/)
- [ ] [Create a new merge request](https://docs.gitlab.com/ee/user/project/merge_requests/creating_merge_requests.html)
- [ ] [Automatically close issues from merge requests](https://docs.gitlab.com/ee/user/project/issues/managing_issues.html#closing-issues-automatically)
- [ ] [Enable merge request approvals](https://docs.gitlab.com/ee/user/project/merge_requests/approvals/)
- [ ] [Set auto-merge](https://docs.gitlab.com/ee/user/project/merge_requests/merge_when_pipeline_succeeds.html)

## Test and Deploy

Use the built-in continuous integration in GitLab.

- [ ] [Get started with GitLab CI/CD](https://docs.gitlab.com/ee/ci/quick_start/index.html)
- [ ] [Analyze your code for known vulnerabilities with Static Application Security Testing (SAST)](https://docs.gitlab.com/ee/user/application_security/sast/)
- [ ] [Deploy to Kubernetes, Amazon EC2, or Amazon ECS using Auto Deploy](https://docs.gitlab.com/ee/topics/autodevops/requirements.html)
- [ ] [Use pull-based deployments for improved Kubernetes management](https://docs.gitlab.com/ee/user/clusters/agent/)
- [ ] [Set up protected environments](https://docs.gitlab.com/ee/ci/environments/protected_environments.html)

***

# Editing this README

When you're ready to make this README your own, just edit this file and use the handy template below (or feel free to structure it however you want - this is just a starting point!). Thanks to [makeareadme.com](https://www.makeareadme.com/) for this template.

## Suggestions for a good README

Every project is different, so consider which of these sections apply to yours. The sections used in the template are suggestions for most open source projects. Also keep in mind that while a README can be too long and detailed, too long is better than too short. If you think your README is too long, consider utilizing another form of documentation rather than cutting out information.

## Name
Choose a self-explaining name for your project.

## Description
Let people know what your project can do specifically. Provide context and add a link to any reference visitors might be unfamiliar with. A list of Features or a Background subsection can also be added here. If there are alternatives to your project, this is a good place to list differentiating factors.

## Badges
On some READMEs, you may see small images that convey metadata, such as whether or not all the tests are passing for the project. You can use Shields to add some to your README. Many services also have instructions for adding a badge.

## Visuals
Depending on what you are making, it can be a good idea to include screenshots or even a video (you'll frequently see GIFs rather than actual videos). Tools like ttygif can help, but check out Asciinema for a more sophisticated method.

## Installation
Within a particular ecosystem, there may be a common way of installing things, such as using Yarn, NuGet, or Homebrew. However, consider the possibility that whoever is reading your README is a novice and would like more guidance. Listing specific steps helps remove ambiguity and gets people to using your project as quickly as possible. If it only runs in a specific context like a particular programming language version or operating system or has dependencies that have to be installed manually, also add a Requirements subsection.

## Usage
Use examples liberally, and show the expected output if you can. It's helpful to have inline the smallest example of usage that you can demonstrate, while providing links to more sophisticated examples if they are too long to reasonably include in the README.

## Support
Tell people where they can go to for help. It can be any combination of an issue tracker, a chat room, an email address, etc.

## Roadmap
If you have ideas for releases in the future, it is a good idea to list them in the README.

## Contributing
State if you are open to contributions and what your requirements are for accepting them.

For people who want to make changes to your project, it's helpful to have some documentation on how to get started. Perhaps there is a script that they should run or some environment variables that they need to set. Make these steps explicit. These instructions could also be useful to your future self.

You can also document commands to lint the code or run tests. These steps help to ensure high code quality and reduce the likelihood that the changes inadvertently break something. Having instructions for running tests is especially helpful if it requires external setup, such as starting a Selenium server for testing in a browser.

## Authors and acknowledgment
Show your appreciation to those who have contributed to the project.

## License
For open source projects, say how it is licensed.

## Project status
If you have run out of energy or time for your project, put a note at the top of the README saying that development has slowed down or stopped completely. Someone may choose to fork your project or volunteer to step in as a maintainer or owner, allowing your project to keep going. You can also make an explicit request for maintainers.
