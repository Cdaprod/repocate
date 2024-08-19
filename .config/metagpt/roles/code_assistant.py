import asyncio
from metagpt.roles.role import Role
from metagpt.schema import Message
from metagpt.actions.write_code import WriteCode
from metagpt.actions.review_code import ReviewCode
from metagpt.actions.optimize_code import OptimizeCode
from metagpt.actions.debug_code import DebugCode
from metagpt.actions.document_code import DocumentCode

class CodeAssistant(Role):
    name = "Code Assistant"
    profile = "A role designed to assist with coding tasks, including writing, reviewing, and optimizing code."

    def __init__(self, **kwargs):
        super().__init__(**kwargs)
        self.set_actions([
            WriteCode(),
            ReviewCode(),
            OptimizeCode(),
            DebugCode(),
            DocumentCode()
        ])

    async def _act(self) -> Message:
        todo = self.rc.todo
        msg = self.get_memories(k=1)[0]
        code_text = await todo.run(msg.content)
        return Message(content=code_text, role=self.profile, cause_by=type(todo))

async def main():
    msg = "Write a function that calculates the product of a list"
    context = Context()
    role = CodeAssistant(context=context)
    result = await role.run(msg)
    print(result)

if __name__ == "__main__":
    asyncio.run(main())
