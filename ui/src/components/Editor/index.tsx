import React from 'react';
import { useQuery } from '@tanstack/react-query';
import { Skeleton, Tabs } from "antd";

import { markdown, markdownLanguage } from '@codemirror/lang-markdown';
import { EditorView, ViewUpdate } from '@codemirror/view';
import { languages } from '@codemirror/language-data';
import CodeMirror from '@uiw/react-codemirror';

import { Markdown } from 'services';

import './style.scss'

export interface IEditorProps {
  value?: string
  onChange?: (value: string) => void
}

const Editor = (props: IEditorProps) => {

  const [activeKey, setActiveKey] = React.useState<string>('write')
  const [content, setContent] = React.useState<string>(props.value || '')

  const {
    isFetching,
    data: html,
  } = useQuery<string>(['markdown.preview', content], () => Markdown.preview(content), {
    enabled: content !== '' && activeKey === 'preview',
    initialData: '',
  })

  const changeHandle = React.useCallback((value: string, viewUpdate: ViewUpdate) => {

    setContent(value);

    if (props.onChange !== undefined) {
      props.onChange(value)
    }
  }, [props])

  React.useEffect(() => {
    setContent(props.value || '');
  }, [props.value])

  return (
    <Tabs
      type="card"
      activeKey={activeKey}
      onChange={(activeKey: string) => setActiveKey(activeKey) }
      items={[
        {
          key: 'write',
          label: `Write`,
          children: <>
            <CodeMirror
              className="code-editor"
              basicSetup={{
                lineNumbers: false,
                foldGutter: false,
                highlightActiveLine: false,
              }}
              extensions={[
                markdown({ base: markdownLanguage, codeLanguages: languages }),
                EditorView.lineWrapping,
              ]}
              value={props.value}
              placeholder="Enter the content here"
              onChange={changeHandle}
            />
          </>,
        },
        {
          key: 'preview',
          label: `Preview`,
          children: <div className="code-editor-preview">
            {
              isFetching
                ? <Skeleton active />
                : <div dangerouslySetInnerHTML={{ __html: html }} />
            }
          </div>,
        },
      ]}
    />
  );
}

export default Editor