// Commitlint configuration for conventional commits
// https://www.conventionalcommits.org/

module.exports = {
  extends: ['@commitlint/config-conventional'],
  rules: {
    // Enforce these commit types (matching Dependabot prefixes)
    'type-enum': [
      2,
      'always',
      [
        'feat',     // New feature
        'fix',      // Bug fix
        'docs',     // Documentation only
        'style',    // Code style (formatting, semicolons, etc.)
        'refactor', // Code change that neither fixes a bug nor adds a feature
        'perf',     // Performance improvement
        'test',     // Adding or fixing tests
        'chore',    // Maintenance tasks
        'ci',       // CI/CD changes
        'deps',     // Dependency updates (used by Dependabot)
        'build',    // Build system changes
        'revert',   // Revert a previous commit
      ],
    ],
    // Allow any scope (component/module name)
    'scope-case': [2, 'always', 'lower-case'],
    // Subject must be lowercase
    'subject-case': [2, 'always', 'lower-case'],
    // No period at end of subject
    'subject-full-stop': [2, 'never', '.'],
    // Reasonable max length
    'header-max-length': [2, 'always', 100],
    // Body must have blank line before it
    'body-leading-blank': [1, 'always'],
    // Footer must have blank line before it
    'footer-leading-blank': [1, 'always'],
  },
};
