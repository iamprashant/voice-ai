import {
  MAX_VAR_KEY_LENGHT,
  VAR_ITEM_TEMPLATE,
  getMaxVarNameLength,
  CONTEXT_PLACEHOLDER_TEXT,
  HISTORY_PLACEHOLDER_TEXT,
  QUERY_PLACEHOLDER_TEXT,
  PRE_PROMPT_PLACEHOLDER_TEXT,
} from '@/configs';

/**
 * Get all the variable from given text
 * @param key
 * @param type
 * @returns
 */
export const getNewVar = (
  key: string,
  type?: string,
): { name: string; type: string; defaultvalue: string } => {
  let varWithDefault = {
    ...VAR_ITEM_TEMPLATE,
    type: 'text',
    name: key.slice(0, getMaxVarNameLength(key)),
    defaultvalue: '',
  };

  if (type) varWithDefault.type = type;
  return varWithDefault;
};

// Enhanced regex to capture more variable patterns
const varRegex =
  /\{\{\s*([a-zA-Z_][a-zA-Z0-9_]*(?:\.[a-zA-Z_][a-zA-Z0-9_]*)*)\s*(?:\|[^}]*)?\s*\}\}/g;

export const getVars = (value: string): string[] => {
  if (!value) return [];

  const variables = new Set<string>();
  let match;

  // Extract variables from {{ }} expressions
  while ((match = varRegex.exec(value)) !== null) {
    const fullMatch = match[0];
    const varName = match[1];

    if (
      [
        CONTEXT_PLACEHOLDER_TEXT,
        HISTORY_PLACEHOLDER_TEXT,
        QUERY_PLACEHOLDER_TEXT,
        PRE_PROMPT_PLACEHOLDER_TEXT,
      ].includes(fullMatch)
    ) {
      continue;
    }

    // Apply length filter
    if (varName.length <= MAX_VAR_KEY_LENGHT) {
      variables.add(varName);
    }
  }

  // Extract variables from {% %} statements
  const stmtRegex =
    /\{%\s*(?:if|elif|for\s+\w+\s+in|set\s+\w+\s*=)\s*([^%]*?)\s*%\}/g;
  while ((match = stmtRegex.exec(value)) !== null) {
    const expression = match[1];

    // Remove string literals first (both single and double quotes)
    const withoutStrings = expression.replace(/(['"])[^'"]*\1/g, '');

    // Extract variable names from the cleaned expression
    const exprVarRegex =
      /\b([a-zA-Z_][a-zA-Z0-9_]*(?:\.[a-zA-Z_][a-zA-Z0-9_]*)*)\b/g;
    let exprMatch;

    while ((exprMatch = exprVarRegex.exec(withoutStrings)) !== null) {
      const varName = exprMatch[1];

      // Skip Jinja2 keywords
      if (!isJinja2Keyword(varName) && varName.length <= MAX_VAR_KEY_LENGHT) {
        variables.add(varName);
      }
    }
  }

  return Array.from(variables);
};

// Helper function to check Jinja2 keywords
function isJinja2Keyword(word: string): boolean {
  const keywords = [
    'and',
    'or',
    'not',
    'is',
    'in',
    'true',
    'false',
    'none',
    'for',
    'if',
    'else',
    'elif',
    'endif',
    'endfor',
    'set',
    'range',
  ];
  return keywords.includes(word.toLowerCase());
}
