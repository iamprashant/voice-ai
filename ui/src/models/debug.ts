// export type Inputs = Record<string, string | number | object>;

// // this is the one that will be passed to api and get from api
// export type TextPrompt = {
//   role: PromptRole;
//   content: string;
// };

// export type TextChatCompletePrompt = {
//   prompt: TextPrompt[];
// };

export enum PromptRole {
  system = 'system',
  user = 'user',
  assistant = 'assistant',
}

export type PromptVariable = {
  name: string;
  type: string;
  defaultvalue: string | null;
};
