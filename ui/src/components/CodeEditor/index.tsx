import React, { CSSProperties, useEffect, useState } from 'react';

import { ViewUpdate } from '@codemirror/view';
import { Extension } from '@codemirror/state';
import CodeMirror, { BasicSetupOptions } from '@uiw/react-codemirror';
import { langs } from '@uiw/codemirror-extensions-langs';

import './style.scss'
import { Skeleton, Tabs } from 'antd';

interface CustomCSSProperties extends CSSProperties {
  '--minHeight': string;
}

export interface IEditorProps {
  value?: string;
  onChange?: (value: string) => void;

  className?: string
  minHeight?: number

  placeholder?: string | HTMLElement;
  basicSetup?: boolean | BasicSetupOptions;
  extensions: Extension[];
  lang?: string;
  autoFocus?: boolean
  preview?: (value: string) => Promise<string>;
}

const CodeEditor = (props: IEditorProps) => {

  const [extensions, setExtensions] = useState<Extension[]>(props.extensions);
  const [style, setStyle] = useState<CustomCSSProperties>()

  const [activeKey, setActiveKey] = React.useState<string>('write')

  const [value, setValue] = useState(props.value || '')
  const [isFetching, setIsFetching] = useState(false)
  const [html, setHTML] = React.useState<string>(props.value || '')

  const changeHandle = React.useCallback((value: string, viewUpdate: ViewUpdate) => {
    setValue(value);
    if (props.onChange !== undefined) {
      props.onChange(value);
    }
  }, [props])

  const tabChangeHandle = (activeKey: string) => {
    setActiveKey(activeKey);

    if (activeKey === 'preview' && props.preview) {
      setIsFetching(true)
      props.preview(value)
        .then((html) => setHTML(html))
        .finally(() => setIsFetching(false))
    }
  }

  useEffect(() => {
    if (props.lang) {
      const lang = props.lang as keyof typeof langs;
      setExtensions([...props.extensions, langs[lang]()]);
    }
  }, [props.extensions, props.lang])

  useEffect(() => {
    if (props.minHeight) {
      setStyle({ "--minHeight": `${props.minHeight}px` })
    }
  }, [props.minHeight])

  const editor = () => (
    <CodeMirror
      className={`code-editor ${props.className}`}
      basicSetup={props.basicSetup}
      extensions={extensions}
      value={props.value}
      placeholder={props.placeholder}
      autoFocus={props.autoFocus}
      onChange={changeHandle}
      style={style}
    />
  )

  if (props.preview === undefined) {
    return editor();
  }

  return (
    <Tabs
      type="card"
      activeKey={activeKey}
      onChange={tabChangeHandle}
      items={[
        {
          key: 'write',
          label: `Write`,
          children: editor(),
        },
        {
          key: 'preview',
          label: `Preview`,
          children: isFetching
            ? <div className="code-editor-preview" style={{ ...style }}><Skeleton active /></div>
            : <div className="code-editor-preview" style={{ ...style }} dangerouslySetInnerHTML={{ __html: html }} />
        },
      ]}
    />
  )
}

CodeEditor.defaultProps = {
  extensions: [
    // EditorView.lineWrapping,
  ],
};

export default CodeEditor