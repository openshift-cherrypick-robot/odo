const lightCodeTheme = require('prism-react-renderer/themes/github');
const darkCodeTheme = require('prism-react-renderer/themes/dracula');
const path = require('path');

/** @type {import('@docusaurus/types').DocusaurusConfig} */
module.exports = {
  title: 'odo',
  tagline: 'Fast iterative Kubernetes and OpenShift development',
  url: 'https://odo.dev',
  baseUrl: '/',
  onBrokenLinks: 'throw',
  onBrokenMarkdownLinks: 'warn',
  favicon: 'img/favicon.ico',
  organizationName: 'redhat-developer', // Usually your GitHub org/user name.
  projectName: 'odo', // Usually your repo name.
  plugins: [
      [
          path.resolve(__dirname, 'docusaurus-odo-plugin-segment'),
        {
          apiKey: 'seYXMF0tyHs5WcPsaNXtSEmQk3FqzTz0',
          options: {
            context: {ip: '0.0.0.0'}
          }
        }
      ]
  ],
  themeConfig: {
    navbar: {
      title: 'odo',
      // logo: {
      //   alt: 'My Site Logo',
      //   src: 'img/logo.svg',
      // },
      items: [
        {
          type: 'doc',
          docId: 'intro',
          position: 'left',
          label: 'Docs',
        },
        {to: '/blog', label: 'Blog', position: 'left'},
        {
          href: 'https://github.com/redhat-developer/odo',
          label: 'GitHub',
          position: 'right',
        },
      ],
    },
    footer: {
      style: 'dark',
      links: [
        {
          title: 'Learn',
          items: [
            {
              label: 'Installation',
              to: 'docs/getting-started/installation'
            },
            {
              label: 'Quickstart',
              to: 'docs/getting-started/quickstart'
            },
          ]
        },
        {
          title: 'Community',
          items: [
            {
              label: 'Slack',
              href: 'https://kubernetes.slack.com/archives/C01D6L2NUAG',
            },
            {
              label: 'Meetings',
              href: 'https://calendar.google.com/calendar/u/0/embed?src=gi0s0v5ukfqkjpnn26p6va3jfc@group.calendar.google.com',
            },
          ],
        },
        {
          title: 'More',
          items: [
            {
              label: 'Blog',
              to: 'blog',
            },
            {
              label: 'GitHub',
              href: 'https://github.com/redhat-developer/odo',
            },
            {
              label: 'Twitter',
              href: 'https://twitter.com/rhdevelopers',
            },
          ],
        },
      ],
      // copyright: `Copyright © ${new Date().getFullYear()} My Project, Inc. Built with Docusaurus.`,
    },
    prism: {
      theme: lightCodeTheme,
      darkTheme: darkCodeTheme,
    },
    algolia: {
      appId: 'BH4D9OD16A',
      apiKey: 'e498f97159ee3094d356b8ed95dd405f',
      indexName: 'odo',
      debug: false
    }
  },
  presets: [
    [
      '@docusaurus/preset-classic',
      {
        docs: {
          sidebarPath: require.resolve('./sidebars.js'),
          // Please change this to your repo.
          editUrl:
            'https://github.com/redhat-developer/odo/edit/main/website/',
        },
        blog: {
          showReadingTime: true,
          // Please change this to your repo.
          editUrl:
            'https://github.com/redhat-developer/odo/edit/main/website/blog/',
          blogSidebarTitle: 'All posts',
          blogSidebarCount: 'ALL',
          postsPerPage: 5,
        },
        theme: {
          customCss: require.resolve('./src/css/custom.css'),
        },
      },
    ],
  ],
};
