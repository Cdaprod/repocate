import asyncio
from metagpt.roles.role import Role
from metagpt.schema import Message
from metagpt.actions.github_actions import FetchRepoData, FetchCommitHistory, FetchIssues, FetchPullRequests

class AdvancedGitHubDataExtractor(Role):
    name = "Advanced GitHub Data Extractor"
    profile = "A role designed to extract detailed data from GitHub repositories."

    def __init__(self, **kwargs):
        super().__init__(**kwargs)
        self.set_actions([
            FetchRepoData(),
            FetchCommitHistory(),
            FetchIssues(),
            FetchPullRequests()
        ])

    async def _act(self) -> Message:
        todo = self.rc.todo
        msg = self.get_memories(k=1)[0]
        data = await todo.run(msg.content)
        return Message(content=data, role=self.profile, cause_by=type(todo))

async def main():
    msg = "Extract detailed data from https://github.com/username/repository"
    context = Context()
    role = AdvancedGitHubDataExtractor(context=context)
    result = await role.run(msg)
    print(result)

if __name__ == "__main__":
    asyncio.run(main())
