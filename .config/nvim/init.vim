" ---- Core Settings ----
set nocompatible
set termguicolors
set nu rnu
set completeopt=menuone,noinsert,noselect
set shortmess+=c
set expandtab
set smartindent
set tabstop=4 softtabstop=4
set cmdheight=2
set updatetime=50
set signcolumn=yes
set clipboard=unnamed,unnamedplus
set shortmess+=I          " Skip intro screen

" Highlight Yank
augroup highlight_yank
  autocmd!
  autocmd TextYankPost * silent! lua require'vim.highlight'.on_yank()
augroup END

" Map leader key to spacebar
let mapleader = " "

" iPhone Specific Keybindings
" Map Ctrl + x to save and quit (write and quit)
nnoremap <C-x> :wq<CR>
" Map Ctrl + c to quit without saving (force quit)
nnoremap <C-c> :q!<CR>

" ---- Plugin Management with vim-plug ----
let vimplug_exists = expand('~/.config/nvim/autoload/plug.vim')
if !filereadable(vimplug_exists)
  if !executable('curl')
    echoerr "You have to install curl or first install vim-plug yourself!"
    execute "q!"
  endif
  echo "Installing Vim-Plug..."
  silent exec "!curl -fLo " . shellescape(vimplug_exists) . " --create-dirs https://raw.githubusercontent.com/junegunn/vim-plug/master/plug.vim"
  autocmd VimEnter * PlugInstall
endif

call plug#begin('~/.config/nvim/plugged')

" Plugins List
Plug 'tpope/vim-sensible'                              " Sensible defaults
Plug 'sainnhe/edge'                                     " Color schemes
Plug 'neovim/nvim-lspconfig'                            " LSP
Plug 'SirVer/ultisnips'                                  " Code snippets
Plug 'honza/vim-snippets'                                " Code snippets
Plug 'nvim-lua/popup.nvim'                               " Fuzzy finder dependencies
Plug 'nvim-lua/plenary.nvim'                             " Fuzzy finder dependencies
Plug 'nvim-telescope/telescope.nvim'                     " Fuzzy finder
Plug 'nvim-treesitter/nvim-treesitter', {'do': ':TSUpdate'} " Syntax highlighting
Plug 'nvim-treesitter/playground'                        " Syntax playground
Plug 'kyazdani42/nvim-web-devicons'                      " File explorer icons
Plug 'kyazdani42/nvim-tree.lua'                          " File explorer
Plug 'glepnir/galaxyline.nvim'                           " Status line
Plug 'nvim-telescope/telescope-dap.nvim'                 " Debugging integration
Plug 'mfussenegger/nvim-dap'                             " Debugging
Plug 'mfussenegger/nvim-dap-python'                      " Python debugging
Plug 'pwntester/octo.nvim'                               " GitHub integration
Plug 'zbirenbaum/copilot.lua'                            " GitHub Copilot
Plug 'zbirenbaum/copilot-cmp'                            " Copilot completion
Plug 'hrsh7th/nvim-cmp'                                  " Nvim CMP (Completions)
Plug 'hrsh7th/cmp-nvim-lsp'                              " LSP completions
Plug 'hrsh7th/cmp-buffer'                                " Buffer completions
Plug 'hrsh7th/cmp-path'                                  " Path completions
Plug 'hrsh7th/cmp-cmdline'                               " Cmdline completions
Plug 'L3MON4D3/LuaSnip'                                  " Snippet engine
Plug 'saadparwaiz1/cmp_luasnip'                          " LuaSnip completions
Plug 'nvim-lualine/lualine.nvim'                          " Lualine status line

call plug#end()

" Automatically install missing plugins on startup
autocmd VimEnter *
  \ if len(filter(values(g:plugs), '!isdirectory(v:val.dir)'))
  \| PlugInstall --sync | q
  \| endif

" ---- Plugin Configuration ----

" Lualine Configuration
lua << EOF
require('lualine').setup {
  options = {
    theme = 'onedark',
    section_separators = {'', ''},
    component_separators = {'', ''}
  },
  sections = {
    lualine_a = {'mode'},
    lualine_b = {'branch', 'diff', 'diagnostics'},
    lualine_c = {'filename'},
    lualine_x = {'encoding', 'fileformat', 'filetype'},
    lualine_y = {'progress'},
    lualine_z = {'location'}
  }
}
EOF

" Treesitter Configuration
lua <<EOF
require'nvim-treesitter.configs'.setup {
  highlight = {
    enable = true
  },
  playground = {
    enable = true,
    disable = {},
    updatetime = 25,
    persist_queries = false
  }
}
EOF

" Copilot Configuration
lua << EOF
require("copilot").setup({
  suggestion = { enabled = false },
  panel = { enabled = false },
})
require("copilot_cmp").setup()
EOF

" Nvim CMP Configuration
lua << EOF
local cmp = require'cmp'
cmp.setup({
  snippet = {
    expand = function(args)
      require('luasnip').lsp_expand(args.body)
    end,
  },
  mapping = cmp.mapping.preset.insert({
    ['<C-b>'] = cmp.mapping.scroll_docs(-4),
    ['<C-f>'] = cmp.mapping.scroll_docs(4),
    ['<C-Space>'] = cmp.mapping.complete(),
    ['<C-e>'] = cmp.mapping.abort(),
    ['<CR>'] = cmp.mapping.confirm({ select = true }),
  }),
  sources = cmp.config.sources({
    { name = 'nvim_lsp' },
    { name = 'buffer' },
  })
})
EOF

" LSP Configuration
lua << EOF
local nvim_lsp = require('lspconfig')
local on_attach = function(client, bufnr)
  require('completion').on_attach()
  local function buf_set_keymap(...) vim.api.nvim_buf_set_keymap(bufnr, ...) end
  local function buf_set_option(...) vim.api.nvim_buf_set_option(bufnr, ...) end
  buf_set_option('omnifunc', 'v:lua.vim.lsp.omnifunc')
  local opts = { noremap=true, silent=true }
  buf_set_keymap('n', 'gD', '<Cmd>lua vim.lsp.buf.declaration()<CR>', opts)
  buf_set_keymap('n', 'gd', '<Cmd>lua vim.lsp.buf.definition()<CR>', opts)
  buf_set_keymap('n', 'K', '<Cmd>lua vim.lsp.buf.hover()<CR>', opts)
  buf_set_keymap('n', 'gi', '<cmd>lua vim.lsp.buf.implementation()<CR>', opts)
  buf_set_keymap('n', '<C-k>', '<cmd>lua vim.lsp.buf.signature_help()<CR>', opts)
  buf_set_keymap('n', '<space>wa', '<cmd>lua vim.lsp.buf.add_workspace_folder()<CR>', opts)
  buf_set_keymap('n', '<space>wr', '<cmd>lua vim.lsp.buf.remove_workspace_folder()<CR>', opts)
  buf_set_keymap('n', '<space>wl', '<cmd>lua print(vim.inspect(vim.lsp.buf.list_workspace_folders()))<CR>', opts)
  buf_set_keymap('n', '<space>D', '<cmd>lua vim.lsp.buf.type_definition()<CR>', opts)
  buf_set_keymap('n', '<space>rn', '<cmd>lua vim.lsp.buf.rename()<CR>', opts)
  buf_set_keymap('n', 'gr', '<cmd>lua vim.lsp.buf.references()<CR>', opts)
  buf_set_keymap('n', '<space>e', '<cmd>lua vim.lsp.diagnostic.show_line_diagnostics()<CR>', opts)
  buf_set_keymap('n', '[d', '<cmd>lua vim.lsp.diagnostic.goto_prev()<CR>', opts)
  buf_set_keymap('n', ']d', '<cmd>lua vim.lsp.diagnostic.goto_next()<CR>', opts)
  buf_set_keymap('n', '<space>q', '<cmd>lua vim.lsp.diagnostic.set_loclist()<CR>', opts)
  if client.resolved_capabilities.document_formatting then
    buf_set_keymap("n", "<space>f", "<cmd>lua vim.lsp.buf.formatting()<CR>", opts)
  elseif client.resolved_capabilities.document_range_formatting then
    buf_set_keymap("n", "<space>f", "<cmd>lua vim.lsp.buf.formatting()<CR>", opts)
  end
  if client.resolved_capabilities.document_highlight then
    require('lspconfig').util.nvim_multiline_command [[
    :hi LspReferenceRead cterm=bold ctermbg=red guibg=LightYellow
    :hi LspReferenceText cterm=bold ctermbg=red guibg=LightYellow
    :hi LspReferenceWrite cterm=bold ctermbg=red guibg=LightYellow
    augroup lsp_document_highlight
      autocmd!
      autocmd CursorHold <buffer> lua vim.lsp.buf.document_highlight()
      autocmd CursorMoved <buffer> lua vim.lsp.buf.clear_references()
    augroup END
    ]]
  end
end

local servers = {'pyright', 'gopls', 'rust_analyzer'}
for _, lsp in ipairs(servers) do
  nvim_lsp[lsp].setup {
    on_attach = on_attach,
  }
end
EOF

" ---- Custom Keybindings ----

" General key mappings
nnoremap <Space>v :e ~/.config/nvim/init.vim<CR>

" Leader Keybindings (Normal mode)
nnoremap <leader>ff :Telescope find_files<CR>          " Find files
nnoremap <leader>fb :Telescope buffers<CR>             " List buffers
nnoremap <leader>gs :Git status<CR>                    " Git status (using fugitive)
nnoremap <leader>gd :Git diff<CR>                      " Git diff (using fugitive)
nnoremap <leader>tt :NvimTreeToggle<CR>                " Toggle file explorer
nnoremap <leader>e :lua vim.diagnostic.open_float()<CR> " Show diagnostics
nnoremap <leader>f :lua vim.lsp.buf.formatting()<CR>   " Format code
nnoremap <leader>dr :lua require'dap'.repl.open()<CR>  " Open DAP repl

" Control Keybindings (All modes)
nnoremap <C-x> :wq<CR>                                 " Save and quit in Normal mode
inoremap <C-x> <Esc>:wq<CR>                            " Save and quit in Insert mode
vnoremap <C-x> <Esc>:wq<CR>                            " Save and quit in Visual mode
snoremap <C-x> <Esc>:wq<CR>                            " Save and quit in Select mode
xnoremap <C-x> <Esc>:wq<CR>                            " Save and quit in Visual mode
onoremap <C-x> <Esc>:wq<CR>                            " Save and quit in Operator mode

nnoremap <C-c> :q!<CR>                                 " Force quit without saving in Normal mode
inoremap <C-c> <Esc>:q!<CR>                            " Force quit without saving in Insert mode
vnoremap <C-c> <Esc>:q!<CR>                            " Force quit without saving in Visual mode
snoremap <C-c> <Esc>:q!<CR>                            " Force quit without saving in Select mode
xnoremap <C-c> <Esc>:q!<CR>                            " Force quit without saving in Visual mode
onoremap <C-c> <Esc>:q!<CR>                            " Force quit without saving in Operator mode

" Window Navigation and Resizing
nnoremap <C-h> <C-w>h                                  " Move to left split
nnoremap <C-j> <C-w>j                                  " Move to below split
nnoremap <C-k> <C-w>k                                  " Move to above split
nnoremap <C-l> <C-w>l                                  " Move to right split

" Completion and Other Actions (Insert mode)
inoremap <C-Space> <cmd>lua require'cmp'.complete()<CR> " Trigger completion

" Alt (Meta) Keybindings (Normal mode)
nnoremap <M-b> :lua require'dap'.toggle_breakpoint()<CR> " Toggle breakpoint
nnoremap <M-d> :lua require'dap.ui.widgets'.hover()<CR>  " Show debug hover
nnoremap <M-w> :q<CR>                                    " Close the current window
nnoremap <M-s> :split<CR>                                 " Split window horizontally

" Debugging configuration
nnoremap <silent> <F5> :lua require'dap'.continue()<CR>
nnoremap <silent> <leader>dd :lua require('dap').continue<CR>
nnoremap <silent> <F10> :lua require'dap'.step_over()<CR>
nnoremap <silent> <F11> :lua require'dap'.step_into()<CR>
nnoremap <silent> <F12> :lua require'dap'.step_out()<CR>
nnoremap <silent> <leader>b :lua require'dap'.toggle_breakpoint()<CR>
nnoremap <silent> <leader>B :lua require'dap'.set_breakpoint(vim.fn.input('Breakpoint condition: '))<CR>
nnoremap <silent> <leader>lp :lua require'dap'.set_breakpoint(nil, nil, vim.fn.input('Log point message: '))<CR>
nnoremap <silent> <leader>dr :lua require'dap'.repl.open()<CR>
nnoremap <silent> <leader>dl :lua require'dap'.repl.run_last()<CR>
nnoremap <silent> <leader>dn :lua require('dap-python').test_method()<CR>
vnoremap <silent> <leader>ds <ESC>:lua require('dap-python').debug_selection()<CR>

" Clean startup no message windows
autocmd VimEnter * silent! redraw!