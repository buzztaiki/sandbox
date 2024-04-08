import yargs from "yargs/yargs";
import { z } from "zod";

class PageUrl {
  readonly baseUrl: URL;
  readonly space: string;
  readonly page: string;

  constructor(baseUrl: URL, space: string, page: string) {
    this.baseUrl = baseUrl;
    this.space = space;
    this.page = page;
  }

  static parse(urlStr: string): PageUrl {
    const url = new URL(urlStr);
    const re = new RegExp(String.raw`/wiki/spaces/(?<space>[^/]+)/pages/(?<page>[^/]+)`);
    const m = url.pathname.match(re);
    if (m?.groups == undefined) {
      throw new Error(`invalid page url: ${urlStr}`);
    }
    return new PageUrl(new URL(url.origin), m.groups.space, m.groups.page);
  }
}

const PageSchema = z.object({
  id: z.string(),
  parentId: z.string(),
  spaceId: z.string(),
  title: z.string(),
  body: z.any(),
  status: z.string(),
  version: z.object({
    number: z.number(),
    message: z.string(),
  }),
});
type Page = z.infer<typeof PageSchema>;

class Client {
  readonly baseUrl: URL;
  readonly email: string;
  readonly apiToken: string;

  constructor(baseUrl: URL, email: string, apiToken: string) {
    this.baseUrl = baseUrl;
    this.email = email;
    this.apiToken = apiToken;
  }

  async api(path: string, init?: FetchRequestInit & { query?: Record<string, string> }): Promise<Response> {
    const url = new URL(this.baseUrl);
    url.search = new URLSearchParams(init?.query).toString();
    url.pathname = path;

    const auth = Buffer.from(`${this.email}:${this.apiToken}`).toString("base64");
    const headers = {
      ...init?.headers,
      Accept: "application/json; charset=utf-8",
      "Content-Type": "application/json; charset=utf-8",
      Authorization: `Basic ${auth}`,
    };
    const res = await fetch(url, { ...init, headers, redirect: "follow" });
    if (!res.ok) {
      throw new Error(`response not ok: ${res.status} ${res.statusText}`);
    }
    return res;
  }

  /**
   * get page by id
   * https://developer.atlassian.com/cloud/confluence/rest/v2/api-group-page/#api-pages-id-put
   */
  async getPage(pageId: string): Promise<Page> {
    const res = await this.api(`/wiki/api/v2/pages/${pageId}`, {
      query: { "body-format": "atlas_doc_format" },
    });
    return PageSchema.parse(await res.json());
  }

  /**
   * move page under parentId
   * https://developer.atlassian.com/cloud/confluence/rest/v2/api-group-page/#api-pages-id-put
   */
  async movePage(page: Page, parentId: string): Promise<Page> {
    const body: Page = {
      ...page,
      parentId,
      version: {
        number: page.version.number + 1,
        message: `move to ${parentId}`,
      },
    };
    const res = await this.api(`/wiki/api/v2/pages/${page.id}`, {
      method: "PUT",
      body: JSON.stringify(body),
    });
    return PageSchema.parse(await res.json());
  }
}

function getEnv(key: string): string {
  const value = process.env[key];
  if (value == undefined) {
    throw new Error(`undefined environment variable: ${key}`);
  }
  return value;
}

async function parseArgs(args: string[]) {
  return await yargs(args)
    .scriptName("move-confluence-pages")
    .command("$0 <page> <parent>", "Move confluence page into parent")
    .positional("page", { type: "string", demandOption: true, describe: "page url" })
    .positional("parent", { type: "string", demandOption: true, desc: "move page into this url" })
    .version(false)
    .help()
    .locale("en").argv;
}

const args = await parseArgs(process.argv.slice(2));
const pageUrl = PageUrl.parse(args.page);
const parentPageUrl = PageUrl.parse(args.parent);
const client = new Client(parentPageUrl.baseUrl, getEnv("ATLASSIAN_USER_EMAIL"), getEnv("ATLASSIAN_API_TOKEN"));
await client.movePage(await client.getPage(pageUrl.page), parentPageUrl.page);
